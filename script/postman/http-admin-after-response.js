const CryptoJS = require('crypto-js');

// =================【設定你的簽章 SALT】（需與後端 config/admin.yaml 的 signature.salt 一致）=================
let SIGNATURE_SALT = "~9U7g2R8zgW&dZ_u";
// =====================================================

// 1. 從變數中取出 Pre-request Script 保存的 Key 和 IV
// (請確保 Pre-request Script 有執行 pm.variables.set("key", ...) 與 pm.variables.set("iv", ...))
const randomKeyPlain = pm.variables.get("key");
const randomIvPlain = pm.variables.get("iv");

// 2. 解析後端回傳的 JSON
const oJson = pm.response.json();

// 【修正：改用 let 宣告，後面才能進行字串置換】
let sC = oJson.c;
let sM = oJson.m;
let sR = oJson.r;
let sT = oJson.t;

// =================【驗證回應簽章 (EncryptionMiddleware)】=================
// 後端拼接順序: sTime, sC, sM, sR, SALT，用「,」串接（跟 request 端用「|」不同）
// 必須在 sC/sM/sR 被下面的 Base64URL 還原、解密覆寫「之前」計算，否則簽章會對不上
let sHeaderTime = pm.response.headers.get("Time") || "";
let sHeaderSignature = pm.response.headers.get("Signature") || "";

let aResponseStrings = [sHeaderTime, sC, sM, sR, SIGNATURE_SALT];
let sResponseStrings = aResponseStrings.join(",");
let sExpectedSignature = CryptoJS.MD5(sResponseStrings).toString(CryptoJS.enc.Hex);

pm.test("回應簽章驗證 (Signature)", function () {
    pm.expect(sExpectedSignature).to.eql(sHeaderSignature);
});

if (sExpectedSignature !== sHeaderSignature) {
    console.error("回應簽章驗證失敗！預期:", sExpectedSignature, "實際 Header:", sHeaderSignature);
} else {
    console.log("回應簽章驗證成功:", sExpectedSignature);
}
// =====================================================

// 準備解密器用的 Key 和 IV
const key = CryptoJS.enc.Utf8.parse(randomKeyPlain);
const iv = CryptoJS.enc.Utf8.parse(randomIvPlain);

// 解密 sC
if (sC) {
    // 【修正：sC 的 Base64URL 還原】
    sC = sC.replace(/-/g, '+').replace(/_/g, '/');
    while (sC.length % 4) { // 👈 修正：改用 sC.length
        sC += '=';
    }

    const decrypted = CryptoJS.AES.decrypt(sC, key, {
        iv: iv,
        mode: CryptoJS.mode.CBC,
        padding: CryptoJS.pad.Pkcs7
    });

    const sCode = decrypted.toString(CryptoJS.enc.Utf8);
    console.log("解密後的後端 Code:", sCode);
}

// 解密 sM
if (sM) {
    // 【修正：sM 的 Base64URL 還原】
    sM = sM.replace(/-/g, '+').replace(/_/g, '/');
    while (sM.length % 4) { // 👈 修正：改用 sM.length
        sM += '=';
    }

    const decrypted = CryptoJS.AES.decrypt(sM, key, {
        iv: iv,
        mode: CryptoJS.mode.CBC,
        padding: CryptoJS.pad.Pkcs7
    });

    const sMessage = decrypted.toString(CryptoJS.enc.Utf8);
    console.log("解密後的後端 Message:", sMessage);
}

// 解密 sR
if (sR) {
    // 【修正：sR 的 Base64URL 還原】
    sR = sR.replace(/-/g, '+').replace(/_/g, '/');
    while (sR.length % 4) { // 👈 修正：改用 sR.length
        sR += '=';
    }

    const decrypted = CryptoJS.AES.decrypt(sR, key, {
        iv: iv,
        mode: CryptoJS.mode.CBC,
        padding: CryptoJS.pad.Pkcs7
    });

    const sResult = decrypted.toString(CryptoJS.enc.Utf8);
    console.log("解密後的後端 Result:", sResult);
}

// 解密 sT
if (sT) {
    // 【修正：sT 的 Base64URL 還原】
    sT = sT.replace(/-/g, '+').replace(/_/g, '/');
    while (sT.length % 4) { // 👈 修正：改用 sT.length
        sT += '=';
    }

    const decrypted = CryptoJS.AES.decrypt(sT, key, {
        iv: iv,
        mode: CryptoJS.mode.CBC,
        padding: CryptoJS.pad.Pkcs7
    });

    const sTotal = decrypted.toString(CryptoJS.enc.Utf8);
    console.log("解密後的後端 Total:", sTotal);
}