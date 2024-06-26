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
import pytest

from paas_wl.bk_app.applications.models import WlApp
from paas_wl.bk_app.processes.controllers import list_processes
from paas_wl.bk_app.processes.entities import Runtime
from paas_wl.bk_app.processes.kres_entities import Instance, Process, Schedule
from paas_wl.infras.resources.generation.version import AppResVerManager

pytestmark = pytest.mark.django_db(databases=["default", "workloads"])


def make_process(wl_app: WlApp, process_type: str) -> Process:
    process = Process(
        app=wl_app,
        name="should-set-by-mapper",
        version=1,
        replicas=1,
        type=process_type,
        schedule=Schedule(cluster_name="", tolerations=[], node_selector={}),
        runtime=Runtime(
            envs={},
            image=process_type,
            command=[],
            args=[],
        ),
    )
    process.name = AppResVerManager(wl_app).curr_version.proc_resources(process).deployment_name
    return process


def test_list_processes(bk_stag_env, wl_app, wl_release, mock_reader):
    mock_reader.set_processes([make_process(wl_app, "web"), make_process(wl_app, "worker")])
    mock_reader.set_instances(
        [
            Instance(app=wl_app, name="web", process_type="web"),
            Instance(app=wl_app, name="worker", process_type="worker"),
        ]
    )

    web_proc = make_process(wl_app, "web")
    web_proc.instances = [Instance(process_type="web", app=wl_app, name="web")]
    worker_proc = make_process(wl_app, "worker")
    worker_proc.instances = [Instance(process_type="worker", app=wl_app, name="worker")]
    assert list_processes(bk_stag_env).processes == [web_proc, worker_proc]


def test_list_processes_with_dirty_release(bk_stag_env, wl_app, wl_dirty_release, mock_reader):
    assert not list_processes(bk_stag_env).processes


def test_list_processes_boundary_case(bk_stag_env, wl_app, wl_release, mock_reader):
    mock_reader.set_processes(
        # worker 没有实例, 不会被忽略
        # beat 未定义在 Procfile, 不会被忽略
        [
            make_process(wl_app, "web"),
            make_process(wl_app, "worker"),
            make_process(wl_app, "beat"),
        ]
    )
    mock_reader.set_instances(
        [Instance(app=wl_app, name="web", process_type="web"), Instance(app=wl_app, name="beat", process_type="beat")]
    )

    web_proc = make_process(wl_app, "web")
    web_proc.instances = [Instance(app=wl_app, name="web", process_type="web")]
    worker_proc = make_process(wl_app, "worker")
    beat_proc = make_process(wl_app, "beat")
    beat_proc.instances = [Instance(app=wl_app, name="beat", process_type="beat")]
    assert list_processes(bk_stag_env).processes == [web_proc, worker_proc, beat_proc]


def test_list_processes_without_release(bk_stag_env, wl_app, wl_release, mock_reader):
    """没有发布过的 WlApp，也能获取进程信息"""
    mock_reader.set_processes([make_process(wl_app, "web")])
    mock_reader.set_instances([Instance(app=wl_app, name="web", process_type="web")])
    wl_release.delete()
    assert len(list_processes(bk_stag_env).processes) == 1
