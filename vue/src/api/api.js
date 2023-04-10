import {get, post} from '@/utils/request'

export function login(data) {
    return post('/login', data)
}

export function getDeviceStatus() {
    return get('/status','' )
}
export function startShell(data) {
    return post('/device/start', data)
}
export function stopShell(data) {
    return post('/device/stop', data)
}