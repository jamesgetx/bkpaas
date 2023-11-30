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
import logging

from six import ensure_text

from paas_wl.bk_app.deploy.app_res.controllers import BkAppHookLogFetcher
from paas_wl.infras.resources.base.exceptions import ReadTargetStatusTimeout
from paas_wl.utils.kubestatus import check_pod_health_status
from paasng.platform.engine.constants import JobStatus
from paasng.platform.engine.exceptions import StepNotInPresetListError
from paasng.platform.engine.models import DeployPhaseTypes
from paasng.platform.engine.utils.output import Style
from paasng.platform.engine.workflow import DeployStep

logger = logging.getLogger(__name__)


def generate_pre_release_hook_name(bkapp_name: str, deploy_id: int) -> str:
    """获取钩子 pod 名称. 需要和 operator 中的保持一致"""
    return f"pre-rel-{bkapp_name}-{deploy_id}"


class PreReleaseDummyExecutor(DeployStep):
    """
    Dummy executor for BkApp pre-release hook

    用于前端正确渲染"执行部署前置命令"步骤, 同时获取日志
    """

    PHASE_TYPE = DeployPhaseTypes.RELEASE

    def start(self, hook_name: str):
        # self._mark_step_start()
        self._perform(hook_name)

    def _mark_step_start(self):
        try:
            step = self.phase.get_step_by_name("执行部署前置命令")
        except StepNotInPresetListError:
            return
        step.mark_procedure_status(JobStatus.PENDING)

    def _perform(self, hook_name: str):
        self.stream.write_message(Style.Warning("Starting pre-release phase"))

        wl_app = self.engine_app.to_wl_obj()
        fetcher = BkAppHookLogFetcher(wl_app)

        try:
            fetcher.wait_for_logs_readiness(wl_app.namespace, hook_name)
        except ReadTargetStatusTimeout as e:
            pod = e.extra_value
            if pod is None:
                self.stream.write_message(
                    Style.Error("Pod is not created normally, please contact the cluster administrator.")
                )
            else:
                health_status = check_pod_health_status(pod)
                self.stream.write_message(Style.Error(health_status.message))
            return

        self.stream.write_title("executing...")

        try:
            for line in fetcher.fetch_logs(wl_app.namespace, hook_name, follow=True):
                self.stream.write_message(ensure_text(line))
        except Exception:
            logger.exception(f"A critical error happened during execute hook({hook_name})")
            self.stream.write_message(Style.Error("fetch logs failed"))
        else:
            self.stream.write_message(Style.Error("pre-release execution succeed."))
