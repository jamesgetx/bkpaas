# -*- coding: utf-8 -*-
"""
TencentBlueKing is pleased to support the open source community by making
蓝鲸智云 - PaaS 平台 (BlueKing - PaaS System) available.
Copyright (C) 2017 THL A29 Limited, a Tencent company. All rights reserved.
Licensed under the MIT License (the "License"); you may not use this file except
in compliance with the License. You may obtain a copy of the License at

    http://opensource.org/licenses/MIT

Unless required by applicable law or agreed to in writing, software distributed under
the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
either express or implied. See the License for the specific language governing permissions and
limitations under the License.

We undertake not to change the open source license (MIT license) applicable
to the current version of the project delivered to anyone in the future.
"""
from pathlib import Path
from unittest import mock

import jinja2
import pytest
from django.conf import settings
from django.utils.translation import gettext_lazy as _

from paasng.monitoring.monitor.alert_rules.ascode.client import AsCodeClient
from paasng.monitoring.monitor.alert_rules.config.constants import DEFAULT_RULE_CONFIGS
from paasng.monitoring.monitor.models import AppAlertRule
from tests.utils.helpers import generate_random_string

random_vhost = generate_random_string()


@pytest.fixture(scope="module", autouse=True)
def mock_import_configs():
    with mock.patch.object(AsCodeClient, "_apply_rule_configs", return_value=None) as mock_method:
        yield mock_method


@pytest.fixture(scope="module", autouse=True)
def mock_get_vhost():
    with mock.patch.dict(
        'paasng.monitoring.monitor.alert_rules.config.metric_label.LABEL_VALUE_QUERY_FUNCS',
        {'vhost': lambda app_code, run_env, module_name: random_vhost},
    ):
        yield


@pytest.fixture
def bk_app_init_rule_configs(bk_app):

    tpl_dir = Path(settings.BASE_DIR) / 'paasng' / 'monitoring' / 'monitor' / 'alert_rules' / 'ascode'
    loader = jinja2.FileSystemLoader([tpl_dir / 'rules_tpl', tpl_dir / 'notice_tpl'])
    j2_env = jinja2.Environment(loader=loader, trim_blocks=True)

    app_code = bk_app.code
    module_scoped_configs = DEFAULT_RULE_CONFIGS['module_scoped']
    notice_group_name = f"[{app_code}] {_('通知组')}"

    default_rules = {
        'high_cpu_usage': {
            'alert_rule_name_format': f'{app_code}-default-{{env}}-high_cpu_usage',
            'template_name': 'high_cpu_usage.yaml.j2',
        },
        'high_mem_usage': {
            'alert_rule_name_format': f'{app_code}-default-{{env}}-high_mem_usage',
            'template_name': 'high_mem_usage.yaml.j2',
        },
        'pod_restart': {
            'alert_rule_name_format': f'{app_code}-default-{{env}}-pod_restart',
            'template_name': 'pod_restart.yaml.j2',
        },
        'oom_killed': {
            'alert_rule_name_format': f'{app_code}-default-{{env}}-oom_killed',
            'template_name': 'oom_killed.yaml.j2',
        },
        'high_rabbitmq_queue_messages': {
            'alert_rule_name_format': f'{app_code}-default-{{env}}-high_rabbitmq_queue_messages',
            'template_name': 'high_rabbitmq_queue_messages.yaml.j2',
        },
    }

    init_rule_configs = {}
    for alert_code, c in default_rules.items():
        for env in ['stag', 'prod']:
            alert_rule_name = c['alert_rule_name_format'].format(env=env)
            init_rule_configs[f"rule/{alert_rule_name}.yaml"] = j2_env.get_template(c['template_name']).render(
                alert_rule_display_name=f"[{app_code}:default:{env}] "
                f"{module_scoped_configs[alert_code]['display_name']}",
                app_code=app_code,
                run_env=env,
                alert_code=alert_code,
                enabled=True,
                metric_labels={'namespace': f'bkapp-{app_code}-{env}', 'vhost': random_vhost},
                namespace=f'bkapp-{app_code}-{env}',
                threshold_expr=module_scoped_configs[alert_code]['threshold_expr'],
                notice_group_name=notice_group_name,
            )

    init_rule_configs['notice/default_notice.yaml'] = j2_env.get_template('notice.yaml.j2').render(
        notice_group_name=f"[{app_code}] {_('通知组')}", receivers=bk_app.get_developers()
    )
    return init_rule_configs


@pytest.fixture
def cpu_usage_alert_rule_obj(bk_app):
    return AppAlertRule.objects.create(
        alert_code='high_cpu_usage',
        display_name='high_cpu_usage',
        enabled=True,
        threshold_expr='>= 0.8',
        receivers=bk_app.get_developers(),
        application=bk_app,
        environment='stag',
        module=bk_app.get_default_module(),
    )
