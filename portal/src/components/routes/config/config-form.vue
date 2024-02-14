<script setup lang="ts">
import { ref } from "vue";
import {
  ElDialog,
  ElForm,
  FormRules,
  ElScrollbar,
  FormInstance,
} from "element-plus";

const LOGSOURCE_TYPE = [
  { id: 1, name: "s3_sqs" },
  { id: 2, name: "Cloudwatch" },
];

const LOGTYPE = ["cloudtrail", "vpc_flow_log", "route53"];

const isShowForm = ref(false);
const itemSelected = ref("");
const formModel = ref<any>({});
const rulesFormRef = ref<FormInstance>();

const formRules = ref<FormRules>({
  queueUrl: {
    required: true,
    message: "Required field",
  },

  // apiTimeout: {
  //   required: true,
  //   message: "Required field",
  // },
  // waitTime: {
  //   required: true,
  //   message: "Required field",
  // },
  // secretKey: {
  //   required: true,
  //   message: "Required field",
  // },
  // accessKey: {
  //   required: true,
  //   message: "Required field",
  // },
  logType: {
    required: true,
    message: "Required field",
  },
  visibilityTimeout: {
    required: true,
    message: "Required field",
  },
  logGroup: {
    required: true,
    message: "Required field",
  },
  // regionName: {
  //   required: true,
  //   message: "Required field",
  // },
  // logStreams: {
  //   required: true,
  //   message: "Required field",
  // },
  // proxy: {
  //   required: true,
  //   message: "Required field",
  // },
});

const showFormDialog = (val?: any) => {
  isShowForm.value = true;
  if (JSON.stringify(val) === "{}") {
    itemSelected.value = LOGSOURCE_TYPE[0].name;
    formModel.value = {};
    return;
  }
  itemSelected.value = val?.protocol;
  formModel.value.id = val?.id;
  formModel.value.queueUrl = val?.properties.queue_url;
  formModel.value.apiTimeout = val?.properties.api_timeout;
  formModel.value.waitTime = val?.properties.wait_time;
  formModel.value.secretKey = val?.properties.secret_access_key;
  formModel.value.accessKey = val?.properties.access_key_id;
  formModel.value.logType = val?.properties.log_type;
  formModel.value.visibilityTimeout = val?.properties.visibility_timeout;
  formModel.value.logGroup = val?.properties.log_group_arn;
  formModel.value.regionName = val?.properties.region_name;
  formModel.value.logStreams = val?.properties.log_streams;
  formModel.value.proxy = val?.properties.proxy;
};

const selectedItem = (val: any) => {
  itemSelected.value = val;
};

defineExpose({
  showFormDialog,
});
// v-if="!formModel.id || "
</script>

<template>
  <el-dialog
    v-model="isShowForm"
    width="650px"
    height="800px"
    center
    :close-on-click-modal="false"
    :destroy-on-close="true"
    :title="formModel?.id ? 'Edit' : 'Add'"
  >
    <template #default>
      <div class="config-form-wrap">
        <div class="config-form">
          <div class="config-form-left">
            <div
              class="type-item"
              :class="itemSelected === 's3_sqs' ? 'type-item-checked' : ''"
              @click="selectedItem('s3_sqs')"
              v-if="
                !formModel.id || (formModel.id && itemSelected === 's3_sqs')
              "
            >
              s3_sqs
            </div>
            <div
              class="type-item"
              :class="itemSelected === 'Cloudwatch' ? 'type-item-checked' : ''"
              @click="selectedItem('Cloudwatch')"
              v-if="
                !formModel.id || (formModel.id && itemSelected === 'Cloudwatch')
              "
            >
              Cloudwatch
            </div>
          </div>
          <div class="config-form-right">
            <template v-if="itemSelected === LOGSOURCE_TYPE[0].name">
              <el-form
                ref="rulesFormRef"
                :model="formModel"
                :rules="formRules"
                label-position="top"
                require-asterisk-position="right"
                :inline="true"
              >
                <el-row>
                  <el-col :span="12">
                    <el-form-item label="Queue url" prop="queueUrl">
                      <el-input v-model="formModel.queueUrl" />
                    </el-form-item>
                    <el-form-item label="Log type" prop="logType">
                      <el-select
                        v-model="formModel.logType"
                        placeholder="Select"
                        style="width: 240px"
                      >
                        <el-option
                          v-for="item in LOGTYPE"
                          :key="item"
                          :label="item"
                          :value="item"
                        />
                      </el-select>
                    </el-form-item>
                    <el-form-item
                      label="Visibility timeout"
                      prop="visibilityTimeout"
                    >
                      <el-input v-model="formModel.visibilityTimeout" />
                    </el-form-item>
                    <el-form-item label="Secret access key" prop="secretKey">
                      <el-input v-model="formModel.secretKey" />
                    </el-form-item>
                  </el-col>
                  <el-col :span="12">
                    <el-form-item label="Access key id" prop="accessKey">
                      <el-input v-model="formModel.accessKey" />
                    </el-form-item>
                    <el-form-item label="Api timeout" prop="apiTimeout">
                      <el-input v-model="formModel.apiTimeout" />
                    </el-form-item>
                    <el-form-item label="Wait time" prop="waitTime">
                      <el-input v-model="formModel.waitTime" />
                    </el-form-item>
                    <el-form-item label="proxy" prop="proxy">
                      <el-input v-model="formModel.proxy" />
                    </el-form-item>
                  </el-col>
                </el-row>
              </el-form>
            </template>

            <template v-if="itemSelected === LOGSOURCE_TYPE[1].name">
              <el-form
                ref="rulesFormRef"
                :model="formModel"
                :rules="formRules"
                label-position="top"
                require-asterisk-position="right"
                :inline="true"
              >
                <el-row>
                  <el-col :span="12">
                    <el-form-item label="Log group arn" prop="logGroup">
                      <el-input v-model="formModel.logGroup" />
                    </el-form-item>
                    <el-form-item label="Log type" prop="logType">
                      <el-select
                        v-model="formModel.logType"
                        placeholder="Select"
                        style="width: 240px"
                      >
                        <el-option
                          v-for="item in LOGTYPE"
                          :key="item"
                          :label="item"
                          :value="item"
                        />
                      </el-select>
                    </el-form-item>
                    <el-form-item label="Region name" prop="regionName">
                      <el-input v-model="formModel.regionName" />
                    </el-form-item>
                    <el-form-item label="Log streams" prop="logStreams">
                      <el-input v-model="formModel.logStreams" />
                    </el-form-item>
                  </el-col>
                  <el-col :span="12">
                    <el-form-item label="Secret access key" prop="secretKey">
                      <el-input v-model="formModel.secretKey" />
                    </el-form-item>
                    <el-form-item label="Access key id" prop="accessKey">
                      <el-input v-model="formModel.accessKey" />
                    </el-form-item>

                    <el-form-item label="proxy" prop="proxy">
                      <el-input v-model="formModel.proxy" /> </el-form-item
                  ></el-col>
                </el-row>
              </el-form>
            </template>
          </div>
        </div>
        <div class="config-form-btn">
          <div class="submit-btn" @click="isShowForm = false">Save</div>
        </div>
      </div>
    </template>
  </el-dialog>
</template>

<style scoped lang="scss">
.config-form {
  display: flex;
}
/* LEFT  */
.config-form-left {
  display: flex;
  flex-direction: column;
  border-right: 1px solid black;
  padding-right: 10px;
  gap: 10px;
}

.type-item {
  width: 100px;
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1px solid black;
  cursor: pointer;
  font-weight: 700;
  border-radius: 4px;
}
.type-item-checked {
  background-color: #41bbf8f8;
  color: white;
  border: 1px solid #41bbf8f8;
}

.config-form-right {
  width: 100%;
  height: 300px;
  padding-left: 30px;
}
.config-form-btn {
  margin-top: 20px;
  display: flex;
  justify-content: end;
}
.submit-btn {
  height: 40px;
  width: 100px;
  display: flex;
  align-items: center;
  justify-content: center;
  // border: 1px solid black;
  cursor: pointer;
  background-color: #41bbf8f8;
  border-radius: 4px;
  color: white;
}
</style>
