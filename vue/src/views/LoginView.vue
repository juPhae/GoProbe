<template>
    <div class="login">
        <h1>Login</h1>
        <form>
            <div class="form-group">
                <label for="username">Username:</label>
                <input type="text" id="username" v-model="username">
            </div>
            <div class="form-group">
                <label for="password">Password:</label>
                <input type="password" id="password" v-model="password">
            </div>
            <button type="submit" @click.prevent="login">Login</button>
        </form>
    </div>
</template>

<script>
import { login } from '@/api/api'
import  { encryptAES }   from '@/utils/crypto'
import { Message } from 'element-ui';

export default {
    data() {
        return {
            username: '',
            password: ''
        }
    },
    methods: {

        login() {
            // 检查用户名和密码是否为空
            if (!this.username || !this.password) {
                Message.error('请输入用户名和密码');
                return;
            }

            const key = '1234567812345678'; // 密钥，需要和后端协商
            const iv = '1234567887654321'; // 偏移量，需要和后端协商

            // 加密用户名和密码
            const encryptedUsername = encryptAES(this.username, key, iv);
            const encryptedPassword = encryptAES(this.password, key, iv);

            // 清空密码
            this.password = '';

            // 发起登录请求
            login({ username: encryptedUsername, password: encryptedPassword }).then(response => {
                if (response.token) {
                    // 存储 token 到 localStorage
                    localStorage.setItem('token', response.token);
                    // 登录成功后跳转到之前的页面
                    const redirect = this.$route.query.redirect || '/device'
                    localStorage.setItem('redirect', redirect)
                    this.$router.push({ path: redirect })
                    // 登录成功，跳转到设备视图页面
                    // this.$router.push('/device');
                } else {
                    // 登录失败，提示错误信息
                    Message.error('登录失败：' + response.message);
                }
            }).catch(error => {
                Message.error('出错了' + error);
                console.log(error);

            });

        },
    }
}
</script>

<style>
.login {
    margin: 0 auto;
    width: 400px;
    padding: 20px;
    border: 1px solid #ccc;
    border-radius: 4px;
    text-align: center;
}

.form-group {
    margin-bottom: 20px;
    text-align: left;
}

label {
    display: block;
    margin-bottom: 10px;
}

input[type="text"],
input[type="password"] {
    display: block;
    width: 100%;
    padding: 10px;
    border: 1px solid #ccc;
    border-radius: 4px;
    box-sizing: border-box;
}

button[type="submit"] {
    background-color: #4CAF50;
    color: white;
    padding: 10px 20px;
    border: none;
    border-radius: 4px;
    cursor: pointer;
}
</style>
