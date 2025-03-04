# -*- coding: utf-8 -*-
# TencentBlueKing is pleased to support the open source community by making
# 蓝鲸智云 - PaaS 平台 (BlueKing - PaaS System) available.
# Copyright (C) 2017 THL A29 Limited, a Tencent company. All rights reserved.
# Licensed under the MIT License (the "License"); you may not use this file except
# in compliance with the License. You may obtain a copy of the License at
#
#     http://opensource.org/licenses/MIT
#
# Unless required by applicable law or agreed to in writing, software distributed under
# the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
# either express or implied. See the License for the specific language governing permissions and
# limitations under the License.
#
# We undertake not to change the open source license (MIT license) applicable
# to the current version of the project delivered to anyone in the future.

from unittest.mock import patch

import pytest
from django.test.utils import override_settings
from django_dynamic_fixture import G

from paas_wl.infras.resources.kube_res.exceptions import AppEntityNotFound
from paas_wl.workloads.networking.ingress.certs import DomainWithCert
from paas_wl.workloads.networking.ingress.constants import AppDomainSource
from paas_wl.workloads.networking.ingress.entities import AutoGenDomain
from paas_wl.workloads.networking.ingress.exceptions import (
    DefaultServiceNameRequired,
    EmptyAppIngressError,
    ValidCertNotFound,
)
from paas_wl.workloads.networking.ingress.kres_entities.ingress import ingress_kmodel
from paas_wl.workloads.networking.ingress.managers.domain import (
    CustomDomainIngressMgr,
    IngressDomainFactory,
    SubdomainAppIngressMgr,
    assign_custom_hosts,
)
from paas_wl.workloads.networking.ingress.models import AppDomain, AppDomainCert, AppDomainSharedCert, Domain

pytestmark = pytest.mark.django_db(databases=["default", "workloads"])


@pytest.mark.auto_create_ns
class TestAssignDomains:
    @pytest.fixture()
    def _foo_shared_cert(self, bk_stag_wl_app):
        AppDomainSharedCert.objects.create(
            tenant_id=bk_stag_wl_app.tenant_id, name="foo", cert_data="", key_data="", auto_match_cns="*.foo.com"
        )

    @pytest.mark.usefixtures("_foo_shared_cert")
    @pytest.mark.parametrize(
        "domains",
        [
            [AutoGenDomain("foo.com"), AutoGenDomain("bar.com")],
            [AutoGenDomain("foo.com"), AutoGenDomain("www.foo.com")],
            [AutoGenDomain("foo.com"), AutoGenDomain("www.foo.com", https_enabled=True)],
        ],
    )
    def test_brand_new_domains(self, bk_stag_wl_app, domains):
        assign_custom_hosts(bk_stag_wl_app, domains, "foo-service")

        ingress_mgr = SubdomainAppIngressMgr(bk_stag_wl_app)
        ingress = ingress_mgr.get()

        assert len(ingress.domains) == len(domains)
        assert {d.host for d in domains} == {d.host for d in ingress.domains}

    @pytest.mark.usefixtures("_foo_shared_cert")
    @pytest.mark.parametrize(
        ("hostnames", "domains_https_enabled"),
        [
            (["www.foo.com", "bar.com"], [True, False]),
            (["www.foo.com", "bar.foo.com"], [True, True]),
        ],
    )
    def test_create_https_domains(self, bk_stag_wl_app, hostnames, domains_https_enabled):
        domains = [AutoGenDomain(hostname, https_enabled=True) for hostname in hostnames]
        assign_custom_hosts(bk_stag_wl_app, domains, default_service_name="foo-service")
        domain_count = len(domains)
        ingress = SubdomainAppIngressMgr(bk_stag_wl_app).get()

        assert [d.tls_enabled for d in ingress.domains] == domains_https_enabled
        assert AppDomain.objects.count() == domain_count

    def test_domain_transfer_partially(self, bk_stag_wl_app, bk_prod_wl_app):
        domains_app1 = [AutoGenDomain("foo.com"), AutoGenDomain("bar.com")]
        assign_custom_hosts(bk_stag_wl_app, domains_app1, "foo-service")

        # Transfer "bar.com" to test_domain_transfer_partially
        domains_app1 = [
            AutoGenDomain("bar.com"),
            AutoGenDomain("app-2.com"),
        ]
        assign_custom_hosts(bk_prod_wl_app, domains_app1, "foo-service")

        ingress = SubdomainAppIngressMgr(bk_stag_wl_app).get()
        assert len(ingress.domains) == 1
        hosts = [d.host for d in ingress.domains]
        assert "foo.com" in hosts
        assert "bar.com" not in hosts

        ingress = SubdomainAppIngressMgr(bk_prod_wl_app).get()
        assert len(ingress.domains) == 2
        hosts = [d.host for d in ingress.domains]
        assert "bar.com" in hosts
        assert "app-2.com" in hosts

    def test_domain_transfer_fully(self, bk_stag_wl_app, bk_prod_wl_app):
        domains = [AutoGenDomain("foo.com")]
        assign_custom_hosts(bk_stag_wl_app, domains, "foo-service")

        # Transfer all domains to bk_prod_wl_app
        assign_custom_hosts(bk_prod_wl_app, domains, "foo-service")

        with pytest.raises(AppEntityNotFound):
            SubdomainAppIngressMgr(bk_stag_wl_app).get()

        ingress = SubdomainAppIngressMgr(bk_prod_wl_app).get()
        assert len(ingress.domains) == 1
        assert [d.host for d in ingress.domains] == ["foo.com"]


@pytest.mark.auto_create_ns
class TestSubdomainAppIngressMgrCommon:
    """Test common interfaces for `SubdomainAppIngressMgr`"""

    @pytest.fixture(autouse=True)
    def _setup_data(self, bk_stag_wl_app):
        AppDomain.objects.create(
            app=bk_stag_wl_app, tenant_id=bk_stag_wl_app.tenant_id, host="bar-2.com", source=AppDomainSource.AUTO_GEN
        )

    def test_sync_no_domains(self, bk_stag_wl_app):
        AppDomain.objects.filter(app=bk_stag_wl_app).delete()
        ingress_mgr = SubdomainAppIngressMgr(bk_stag_wl_app)
        with pytest.raises(EmptyAppIngressError):
            ingress_mgr.sync(default_service_name="foo")

    def test_sync_creation_with_no_default_server_name(self, bk_stag_wl_app):
        ingress_mgr = SubdomainAppIngressMgr(bk_stag_wl_app)
        with pytest.raises(DefaultServiceNameRequired):
            ingress_mgr.sync()

    def test_sync_creation(self, bk_stag_wl_app):
        ingresses = ingress_kmodel.list_by_app(bk_stag_wl_app)
        assert len(ingresses) == 0

        ingress_mgr = SubdomainAppIngressMgr(bk_stag_wl_app)
        ingress_mgr.sync(default_service_name="foo")
        ingresses = ingress_kmodel.list_by_app(bk_stag_wl_app)

        assert len(ingresses) == 1
        ingress = ingresses[0]
        assert len(ingress.domains) > 0

    def test_sync_update(self, bk_stag_wl_app):
        ingress_mgr = SubdomainAppIngressMgr(bk_stag_wl_app)
        ingress_mgr.sync(default_service_name="foo")
        ingress_name = ingress_mgr.ingress_name
        assert len(ingress_kmodel.get(bk_stag_wl_app, ingress_name).domains) == 1

        # Add an extra domain
        config = bk_stag_wl_app.latest_config
        config.domain = "bar.com"
        config.save()
        ingress_mgr.sync()

        assert len(ingress_kmodel.get(bk_stag_wl_app, ingress_name).domains) == 2

    def test_delete_non_existed(self, bk_stag_wl_app):
        ingress_mgr = SubdomainAppIngressMgr(bk_stag_wl_app)
        ingress_mgr.delete()

    def test_integrated(self, bk_stag_wl_app):
        ingress_mgr = SubdomainAppIngressMgr(bk_stag_wl_app)
        ingress_mgr.sync(default_service_name="foo")
        assert len(ingress_kmodel.list_by_app(bk_stag_wl_app)) == 1

        ingress_mgr.delete()
        assert len(ingress_kmodel.list_by_app(bk_stag_wl_app)) == 0

    def test_update_target(self, bk_stag_wl_app):
        ingress_mgr = SubdomainAppIngressMgr(bk_stag_wl_app)
        ingress_mgr.sync(default_service_name="foo")
        ingress_mgr.update_target("foo-service", "foo-port")

        ingress = ingress_kmodel.get(bk_stag_wl_app, ingress_mgr.ingress_name)
        assert ingress.service_name == "foo-service"
        assert ingress.service_port_name == "foo-port"

    @pytest.mark.parametrize(
        ("rewrite_to_root", "expected_ret"),
        [
            (True, True),
            (False, False),
        ],
    )
    def test_rewrite_ingress_path_to_root(self, rewrite_to_root, expected_ret, bk_stag_wl_app):
        with patch.object(SubdomainAppIngressMgr, "rewrite_ingress_path_to_root", new=rewrite_to_root):
            SubdomainAppIngressMgr(bk_stag_wl_app).sync(default_service_name="foo")
            assert ingress_kmodel.list_by_app(bk_stag_wl_app)[0].rewrite_to_root is expected_ret


class TestSubdomainAppIngressMgr:
    def test_list_desired_domains(self, bk_stag_wl_app):
        ingress_mgr = SubdomainAppIngressMgr(bk_stag_wl_app)
        domains = ingress_mgr.list_desired_domains()
        assert len(domains) == 0

    def test_list_desired_domains_with_extra(self, bk_stag_wl_app):
        config = bk_stag_wl_app.latest_config
        config.domain = "bar.com"
        config.save()

        ingress_mgr = SubdomainAppIngressMgr(bk_stag_wl_app)
        assert len(ingress_mgr.list_desired_domains()) == 1

        AppDomain.objects.create(
            app=bk_stag_wl_app, tenant_id=bk_stag_wl_app.tenant_id, host="bar-2.com", source=AppDomainSource.AUTO_GEN
        )
        assert len(ingress_mgr.list_desired_domains()) == 2

    def test_list_desired_domains_with_wrong_source(self, bk_stag_wl_app):
        ingress_mgr = SubdomainAppIngressMgr(bk_stag_wl_app)
        assert len(ingress_mgr.list_desired_domains()) == 0
        # SubDomain ingress should only include domains when their source is "CUSTOM"
        AppDomain.objects.create(
            app=bk_stag_wl_app, tenant_id=bk_stag_wl_app.tenant_id, host="foo.com", source=AppDomainSource.INDEPENDENT
        )
        assert len(ingress_mgr.list_desired_domains()) == 0


@pytest.mark.mock_get_structured_app
@pytest.mark.auto_create_ns
class TestCustomDomainIngressMgr:
    @pytest.mark.parametrize(
        ("path_prefix", "expected_path_prefixes", "customized_ingress_name"),
        [
            ("/", ["/"], False),
            ("/foo/", ["/foo/"], True),
        ],
    )
    def test_create(self, path_prefix, expected_path_prefixes, customized_ingress_name, bk_stag_env, bk_stag_wl_app):
        domain = G(
            Domain,
            name="foo.example.com",
            path_prefix=path_prefix,
            module_id=bk_stag_env.module.id,
            environment_id=bk_stag_env.id,
        )
        mgr = CustomDomainIngressMgr(domain)

        mgr.sync(default_service_name=bk_stag_wl_app.name)
        obj = ingress_kmodel.get(bk_stag_wl_app, mgr.make_ingress_name())

        if customized_ingress_name:
            assert obj.name == f"custom-foo.example.com-{domain.id}"
        else:
            assert obj.name == "custom-foo.example.com"
        assert obj.domains[0].path_prefix_list == expected_path_prefixes
        assert obj.service_name == bk_stag_wl_app.name
        mgr.delete()

    def test_normal_delete(self, bk_stag_env, bk_stag_wl_app):
        domain = G(
            Domain,
            name="foo.example.com",
            module_id=bk_stag_env.module.id,
            environment_id=bk_stag_env.id,
        )
        mgr = CustomDomainIngressMgr(domain)
        mgr.sync(default_service_name=bk_stag_wl_app.name)
        _ = ingress_kmodel.get(bk_stag_wl_app, mgr.make_ingress_name())
        mgr.delete()
        with pytest.raises(AppEntityNotFound):
            ingress_kmodel.get(bk_stag_wl_app, mgr.make_ingress_name())

    def test_get_ingress_class(self, bk_stag_env, bk_stag_wl_app):
        domain = G(
            Domain,
            name="foo.example.com",
            module_id=bk_stag_env.module.id,
            environment_id=bk_stag_env.id,
        )
        mgr = CustomDomainIngressMgr(domain)
        assert mgr.get_annotations().get("kubernetes.io/ingress.class") is None

        with override_settings(APP_INGRESS_CLASS="test-ingress"):
            mgr = CustomDomainIngressMgr(domain)
            assert mgr.get_annotations()["kubernetes.io/ingress.class"] == "test-ingress"

        with override_settings(CUSTOM_DOMAIN_INGRESS_CLASS="test-custom-ingress"):
            mgr = CustomDomainIngressMgr(domain)
            assert mgr.get_annotations()["kubernetes.io/ingress.class"] == "test-custom-ingress"


@pytest.mark.auto_create_ns
class TestIntegratedDomains:
    """Test cases for some combined situations"""

    def test_assign_custom_hosts_affects_no_independent_domains(self, bk_stag_wl_app):
        AppDomain.objects.create(
            app=bk_stag_wl_app,
            tenant_id=bk_stag_wl_app.tenant_id,
            host="foo-independent.com",
            source=AppDomainSource.INDEPENDENT,
        )
        assert AppDomain.objects.filter(host="foo-independent.com").exists()
        assign_custom_hosts(bk_stag_wl_app, [AutoGenDomain("foo.com")], "foo-service")

        # Calling assign_custom_hosts should not remove AppDomain objects with source other than "CUSTOM"
        assert AppDomain.objects.filter(host="foo-independent.com").exists()


class TestIngressDomainFactory:
    def test_make_ingress_domain_with_http(self, bk_stag_wl_app):
        domain_with_cert = DomainWithCert(host="example.com", path_prefix="", https_enabled=False)
        factory = IngressDomainFactory(bk_stag_wl_app)
        domain = factory.create(domain_with_cert)
        assert domain.host == domain_with_cert.host
        assert domain.tls_enabled == domain_with_cert.https_enabled
        assert domain.tls_secret_name == ""

    def test_https_cert_not_found(self, bk_stag_wl_app):
        domain_with_cert = DomainWithCert(host="example.com", path_prefix="", https_enabled=True)
        factory = IngressDomainFactory(bk_stag_wl_app)
        with pytest.raises(ValidCertNotFound):
            factory.create(domain_with_cert)

    def test_https_cert_not_found_no_exception(self, bk_stag_wl_app):
        domain_with_cert = DomainWithCert(host="example.com", path_prefix="", https_enabled=True)
        factory = IngressDomainFactory(bk_stag_wl_app)
        domain = factory.create(domain_with_cert, raise_on_no_cert=False)
        assert domain.tls_enabled is False

    @pytest.mark.auto_create_ns
    @pytest.mark.parametrize(
        ("cert_type", "expected_cert_name"),
        [(AppDomainCert, "eng-normal-test"), (AppDomainSharedCert, "eng-shared-test")],
    )
    def test_https_cert_created(self, bk_stag_wl_app, cert_type, expected_cert_name):
        cert = G(cert_type, name="test")
        domain_with_cert = DomainWithCert(host="example.com", path_prefix="", https_enabled=True, cert=cert)
        factory = IngressDomainFactory(bk_stag_wl_app)

        domain = factory.create(domain_with_cert)
        assert domain.host == domain_with_cert.host
        assert domain.tls_enabled == domain_with_cert.https_enabled
        assert domain.tls_secret_name == expected_cert_name
