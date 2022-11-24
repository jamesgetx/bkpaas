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
from django.urls import re_path

from paas_wl.utils import text

from . import views

urlpatterns = [
    re_path(
        f'^regions/{text.PVAR_REGION}/apps/{text.PVAR_NAME}/rcstate_binding/$',
        views.RCStateBindingsViewSet.as_view({'post': 'create', 'delete': 'destroy'}),
    ),
    re_path(
        # use 'services/' prefix to be backward-compatible
        f'^(services/)?regions/{text.PVAR_REGION}/clusters/{text.PVAR_CLUSTER_NAME}/egress_info/$',
        views.ClusterEgressViewSet.as_view({'get': 'retrieve'}),
    ),
]
