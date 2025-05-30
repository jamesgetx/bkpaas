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

# Generated by Django 2.2.17 on 2021-03-03 03:23

from django.db import migrations
from paasng.core.core.storages.sqlalchemy import console_db
from paasng.core.core.storages.utils import DummyDB
from paasng.accessories.publish.sync_market.managers import AppTagManger


def forwards_func(apps, schema_editor):
    if not isinstance(console_db, DummyDB):
        with console_db.session_scope() as session:
            AppTagManger(session).sync_tags_from_console()
    else:
        Tag = apps.get_model("market", "Tag")

        def get_tag_objs():
            for i, n in enumerate(["运维工具", "监控告警", "配置管理", "开发工具", "企业IT", "办公应用", "其它"]):
                yield Tag(id=i + 1, name=n)

        Tag.objects.bulk_create(list(get_tag_objs()))


def reverse_func(apps, schema_editor):
    Tag = apps.get_model("market", "Tag")
    # 最前置的 Tag 初始化脚本，可以安全清除
    Tag.objects.all().delete()


class Migration(migrations.Migration):

    dependencies = [
        ('sync_market', '0001_initial'),
    ]

    operations = [migrations.RunPython(forwards_func, reverse_func)]
