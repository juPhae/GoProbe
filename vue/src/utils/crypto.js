import CryptoJS from 'crypto-js'

function encryptAES(data, key, iv) {
  const encrypted = CryptoJS.AES.encrypt(
    CryptoJS.enc.Utf8.parse(data),
    CryptoJS.enc.Utf8.parse(key),
    {
      iv: CryptoJS.enc.Utf8.parse(iv),
      mode: CryptoJS.mode.CFB,
      padding: CryptoJS.pad.ZeroPadding,
    }
  )

  return encrypted.toString()
}

function decryptAES(data, key, iv) {
  const decrypted = CryptoJS.AES.decrypt(data, CryptoJS.enc.Utf8.parse(key), {
    iv: CryptoJS.enc.Utf8.parse(iv),
    mode: CryptoJS.mode.CFB,
    padding: CryptoJS.pad.ZeroPadding,
  })

  return CryptoJS.enc.Utf8.stringify(decrypted).toString()
}
 // 导出
 export { encryptAES, decryptAES };
