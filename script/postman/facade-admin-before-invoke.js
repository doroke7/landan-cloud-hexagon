
const forge = pm.require("npm:node-forge");
const CryptoJS = require('crypto-js');

// =================【對齊 constants/facade.ts FACADE.ADMIN】=================
const VER = "1.0.1";
const VERSION = "2026-0304";
const SALT = "~9U7g2R8zgW&dZ_u";
// =====================================================

// =================【對齊 constants/http.ts HTTP.ADMIN.PUBLIC_KEY，facade.ts 目前共用同一把 RSA 公鑰】=================
const rsaPublicKeyPem = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAqBbztPHxvEsnb5BcMTjv
693XqRMYte+ORJvrgc0RsdHlC4W4lSqvjnH2JcsTGgtSqmmvsvoPQ9dKGs6+OD3E
zIClDiEr4n7QRFFjYKP3IBqhkR5a5wZdiOCoYCx2dOKjBTkLgIzMO145nITHR0za
Yv7k22eNdIzlLVat1Oq1DlWCHWBEQHUUm/OhiBSHnRb2DXiMa+vBvHHrZBIcDb0+
TRD14zLArY5ijKWkzTLGzr4IDi3TcwDz6xEkLm4grzi/KEYtjAweVTClqm19vYAk
SDe+BtVYNxODv3yQSSIrDEzeCnbimIBCBfwxL65YrbIAUx7YqVbtNry56C4MI95h
rQIDAQAB
-----END PUBLIC KEY-----`;


const ENCRYPTION_FIELDS_MAP = {
    "pb.facade.admin.authentication.Authenticator.SignIn": ["password"],
};

// 對齊 constants/aes.ts AES.CHARSET
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^*()_+-={}|[]";
const charsetLength = charset.length;

function randomString(length) {
    let result = "";
    for (let i = 0; i < length; i++) {
        const randomValue = CryptoJS.lib.WordArray.random(4).words[0];
        const randomIndex = Math.abs(randomValue) % charsetLength;
        result += charset.charAt(randomIndex);
    }
    return result;
}

// base64 轉 base64url，對齊 aesHelper.encrypt / rsaHelper.encrypt 內建的轉換
function toBase64Url(base64Str) {
    return base64Str.replace(/\+/g, '-').replace(/\//g, '_').replace(/=/g, '');
}

// 1. 依「目前呼叫的 RPC」查表，找出這支 RPC 有哪些欄位要加密，
//    對齊 facade.ts 的 oReq.method.input.fields.filter((oField) => getOption(oField, encrypted))
const sRoute = pm.request.methodPath;
const aEncryptedFields = ENCRYPTION_FIELDS_MAP[sRoute] || [];



if (aEncryptedFields.length === 0) {
    console.warn("路由 " + sRoute + " 不在 ENCRYPTION_FIELDS_MAP 裡，本次不加密任何欄位");
}

// 2. 隨機產生本次請求專用的 AES Key / IV，對齊 aesHelper.randomString(16)
const randomKeyPlain = randomString(16);
const randomIvPlain = randomString(16);

const key = CryptoJS.enc.Utf8.parse(randomKeyPlain);
const iv = CryptoJS.enc.Utf8.parse(randomIvPlain);

// pm.request.messages 是 PropertyList，.idx(0).data 底下只有 content 是純文字，
// 要自己 JSON.parse 才能拿到 Message 分頁裡打的 { "name": ..., "password_raw": ... }
let oMessage = {};
try {
    let sContent = pm.request.messages.idx(0).data.content;
    oMessage = JSON.parse(sContent);
} catch (e) {
    console.error("讀取 Message 失敗：", e.message);
}

let aEncryptedFieldNames = [];

aEncryptedFields.forEach((sFieldName) => {
    let sPlainValue = oMessage[sFieldName + "_raw"] || "";

    if (!sPlainValue) {
        return;
    }

    const encrypted = CryptoJS.AES.encrypt(sPlainValue, key, {
        iv: iv,
        mode: CryptoJS.mode.CBC,
        padding: CryptoJS.pad.Pkcs7
    });
    const sEncryptedValue = toBase64Url(encrypted.ciphertext.toString(CryptoJS.enc.Base64));

    pm.variables.set(sFieldName, sEncryptedValue);
    console.log("設定 " + sFieldName + " -> " + sEncryptedValue);
    aEncryptedFieldNames.push(sFieldName);
});

// 4. 用 RSA 公鑰加密本次的 AES Key/IV，對齊 rsaHelper.encrypt(sKeys, CONSTANTS.HTTP.ADMIN.PUBLIC_KEY)
const sKeys = JSON.stringify({ key: randomKeyPlain, iv: randomIvPlain });

let sK = "RSA_ENCRYPTION_FAILED";
try {
    const publicKey = forge.pki.publicKeyFromPem(rsaPublicKeyPem);
    const bytesToEncrypt = forge.util.encodeUtf8(sKeys);
    const rsaEncryptedBytes = publicKey.encrypt(bytesToEncrypt, 'RSAES-PKCS1-V1_5');

    sK = toBase64Url(forge.util.encode64(rsaEncryptedBytes));
} catch (e) {
    console.error("RSA 加密失敗：", e.message);
}

const sTime = Math.floor(Date.now() / 1000).toString();

const sStrings = [VER, VERSION, sK, sTime, SALT].join("|");
const sSignature = CryptoJS.MD5(sStrings).toString(CryptoJS.enc.Hex);

// POSTMAN 的 gPRC 沒辦法 修改 pm.metadata
// 改用環境變數思維
pm.variables.set("VER", VER);
pm.variables.set("VERSION", VERSION);
pm.variables.set("SIGNATURE", sSignature);
pm.variables.set("TIME", sTime);
pm.variables.set("K", sK);


if (aEncryptedFieldNames.length > 0) {
    console.log("已加密並寫入變數（Message 內請用 {{欄位名}} 引用）:", aEncryptedFieldNames.join(", "));
}
