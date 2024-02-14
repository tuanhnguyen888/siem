<script setup lang="ts">
import { ref } from "vue";
import {
  ElForm,
  ElButton,
  ElNotification,
  type FormInstance,
  type FormRules,
} from "element-plus";
import router from "../../common/router";
import store from "../../common/store";

import { userList } from "../demo-data";
import { IUser } from "../../common/model";

const showA = () => {
  ElNotification({
    type: "success",
    title: "Thành Công",
    message: "Đăng nhập thành công",
    duration: 3000,
  });
};
const rulesFormRef = ref<FormInstance>();

const formRules = ref<FormRules>({
  username: [
    {
      required: true,
      message: "Required field",
    },
  ],

  password: [
    {
      required: true,
      message: "Required field",
    },
  ],
});

const formModel = ref<IUser>({});

const resetForm = (elForm: FormInstance | undefined) => {
  if (!elForm) return;

  elForm.resetFields();
};

const login = async () => {
  const data = {
    username: formModel.value.username,
    password: formModel.value.password,
  };

  const isExistUser = userList.filter(
    (item) => item.username === formModel.value.username
  );
  if (!isExistUser.length) {
    return ElNotification({
      type: "error",
      title: "Error",
      message: "Not exist user",
      duration: 1000,
    });
  }
  if (isExistUser[0].password !== formModel.value.password) {
    return ElNotification({
      type: "error",
      title: "Error",
      message: "Wrong password",
      duration: 1000,
    });
  }

  store.dispatch("getLogin");
  ElNotification({
    type: "success",
    title: "Successful",
    message: "Login succesful",
    duration: 1000,
  });
};

const handleSubmit = (elForm: FormInstance | undefined) => {
  if (!elForm) return;

  elForm.validate((valid) => {
    if (valid) {
      login();
    } else {
      return false;
    }
  });
};
</script>

<template>
  <div class="user-login-form" @keydown.enter="handleSubmit(rulesFormRef)">
    <el-form
      ref="rulesFormRef"
      :rules="formRules"
      :model="formModel"
      label-position="top"
      size="large"
      require-asterisk-position="right"
    >
      <el-form-item label="Username" prop="username">
        <el-input v-model="formModel.username" placeholder="Username" />
      </el-form-item>
      <el-form-item label="Password" prop="password">
        <el-input
          v-model="formModel.password"
          type="password"
          placeholder="Nhập mật khẩu"
        />
      </el-form-item>
    </el-form>
    <div class="login-form-btn">
      <el-button
        round
        @click="handleSubmit(rulesFormRef)"
        color="#49b6fb"
        plain
        size="large"
        class="admin-auth-button"
        center
        >Đăng Nhập</el-button
      >
      <el-button
        round
        @click="resetForm(rulesFormRef)"
        color="#49b6fb"
        plain
        size="large"
        class="admin-auth-button"
        center
        >Nhập Lại</el-button
      >
    </div>
    <div class="login-form-footer cursor-pointer" @click="showA">
      Quên mật khẩu
    </div>
  </div>
</template>

<style scoped>
.login-form-btn {
  display: flex;
  justify-content: center;
  color: #ff8c00;
}
.login-form-footer {
  padding-top: 10px;
  display: flex;
  justify-content: center;
}
.login-form-footer:hover {
  text-decoration: underline;
}
</style>
<style lang="scss">
.user-login-form .el-form .el-form-item {
  .el-form-item__label {
    color: rgb(0, 0, 0);
  }
  .el-form-item__label:after {
    color: rgb(255, 0, 0) !important;
  }
  .el-form-item__content {
    .el-form-item__error {
      color: rgb(255, 0, 0);
    }
  }
}
</style>
