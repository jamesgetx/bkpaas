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

from django.urls import path

from paasng.utils.basic import re_path

from . import views

# FIXME: phalanx 方案已经下线, 删除废弃 url
urlpatterns = [
    re_path(
        r"api/monitor/applications/(?P<code>[^/]+)/record/query/$",
        views.EventRecordView.as_view({"post": "query"}),
    ),
    re_path(
        r"api/monitor/applications/(?P<code>[^/]+)/record/(?P<record>[^/]+)/$",
        views.EventRecordDetailsView.as_view({"get": "get"}),
    ),
    re_path(
        r"api/monitor/applications/(?P<code>[^/]+)/record_metrics/(?P<record>[^/]+)/$",
        views.EventRecordMetricsView.as_view({"get": "get"}),
    ),
    re_path(r"api/monitor/applications/(?P<code>[^/]+)/genre/$", views.EventGenreView.as_view({"get": "list"})),
    re_path(
        r"api/monitor/record/applications/summary/$",
        views.EventRecordView.as_view({"get": "app_summary"}),
    ),
    path(
        "api/monitor/applications/<slug:code>/modules/<slug:module_name>/alert_rules/",
        views.AlertRulesView.as_view({"get": "list"}),
    ),
    # TODO init_alert_rules will be deprecated
    path(
        "api/monitor/applications/<slug:code>/alert_rules/init/",
        views.AlertRulesView.as_view({"post": "init_alert_rules"}),
    ),
    path("api/monitor/supported_alert_rules/", views.AlertRulesView.as_view({"get": "list_supported_alert_rules"})),
    path("api/monitor/applications/<slug:code>/alerts/", views.ListAlertsView.as_view({"post": "list"})),
    path("api/monitor/user/alerts/", views.ListAlertsView.as_view({"post": "list_alerts_by_user"})),
    path(
        "api/monitor/applications/<slug:code>/alarm_strategies/",
        views.ListAlarmStrategiesView.as_view({"post": "list"}),
    ),
    path(
        "api/monitor/applications/<slug:code>/dashboard_info/",
        views.GetDashboardInfoView.as_view({"get": "get_dashboard_url"}),
    ),
    path(
        "api/monitor/applications/<slug:code>/builtin_dashboards/",
        views.GetDashboardInfoView.as_view({"get": "list_builtin_dashboards"}),
        name="api.modules.monitor.builtin_dashboards",
    ),
]
