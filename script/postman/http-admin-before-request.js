// 1. 直接從 NPM 載入 node-forge 套件
const forge = pm.require("npm:node-forge");
const CryptoJS = require('crypto-js');

// =================【設定你的簽章 SALT】=================
const SIGNATURE_SALT = "~9U7g2R8zgW&dZ_u"; 
// =====================================================

// 2. 定義允許的自訂字元集與隨機生成函式
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^*()_+-={}|[]";
const charsetLength = charset.length;

function generateRandomString(length) {
    let result = "";
    for (let i = 0; i < length; i++) {
        const randomValue = CryptoJS.lib.WordArray.random(4).words[0];
        const randomIndex = Math.abs(randomValue) % charsetLength;
        result += charset.charAt(randomIndex);
    }
    return result;
}

// 將標準 Base64 轉換為 Base64URL 格式
const toBase64Url = (base64Str) => base64Str
    .replace(/\+/g, '-')
    .replace(/\//g, '_')
    .replace(/=/g, '');

// 3. 隨機產生 16 個字元的 AES Key 與 IV
const randomKeyPlain = generateRandomString(16);
const randomIvPlain = generateRandomString(16);

let oKeys = {
    "key": randomKeyPlain,
    "iv": randomIvPlain,
};
let sKeys = JSON.stringify(oKeys);

// 保存給 Tests (Post-script) 解密使用
pm.variables.set("key", randomKeyPlain);
pm.variables.set("iv", randomIvPlain);

const key = CryptoJS.enc.Utf8.parse(randomKeyPlain);
const iv = CryptoJS.enc.Utf8.parse(randomIvPlain);

// 4. 從 Postman 的 Body 中取出原始 JSON 數據
// 如果打的 body 外面多包了一層 "value"（例如 {"value": {...}}），先拆掉，
// 避免後端 Flatten(..., "value", oValue) 再包一次變成 value.value.xxx。
let jsonString = "";
if (pm.request.body && pm.request.body.raw) {
    try {
        const oRawBody = JSON.parse(pm.request.body.raw);
        const oUnwrapped = (oRawBody && typeof oRawBody === "object" && "value" in oRawBody)
            ? oRawBody.value
            : oRawBody;
        jsonString = JSON.stringify(oUnwrapped);
    } catch (oErr) {
        jsonString = pm.request.body.raw;
    }
} else {
    jsonString = JSON.stringify({ error: "Body was empty" });
}

// 5. 進行 AES-128-CBC 加密
const encrypted = CryptoJS.AES.encrypt(jsonString, key, {
    iv: iv,
    mode: CryptoJS.mode.CBC,
    padding: CryptoJS.pad.Pkcs7
});
const sV = toBase64Url(encrypted.toString());

// =================【RSA 加密 sKeys (得到 sK)】=================
const rsaPublicKeyPem = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAqBbztPHxvEsnb5BcMTjv
693XqRMYte+ORJvrgc0RsdHlC4W4lSqvjnH2JcsTGgtSqmmvsvoPQ9dKGs6+OD3E
zIClDiEr4n7QRFFjYKP3IBqhkR5a5wZdiOCoYCx2dOKjBTkLgIzMO145nITHR0za
Yv7k22eNdIzlLVat1Oq1DlWCHWBEQHUUm/OhiBSHnRb2DXiMa+vBvHHrZBIcDb0+
TRD14zLArY5ijKWkzTLGzr4IDi3TcwDz6xEkLm4grzi/KEYtjAweVTClqm19vYAk
SDe+BtVYNxODv3yQSSIrDEzeCnbimIBCBfwxL65YrbIAUx7YqVbtNry56C4MI95h
rQIDAQAB
-----END PUBLIC KEY-----`;

let sK = "RSA_ENCRYPTION_FAILED";
try {
    const publicKey = forge.pki.publicKeyFromPem(rsaPublicKeyPem);
    const bytesToEncrypt = forge.util.encodeUtf8(sKeys);
    const rsaEncryptedBytes = publicKey.encrypt(bytesToEncrypt, 'RSAES-PKCS1-V1_5');
    
    sK = toBase64Url(forge.util.encode64(rsaEncryptedBytes));
    pm.variables.set("K", sK);
} catch (e) {
    console.error("RSA 加密失敗：", e.message);
}

pm.variables.set("v", sV);

// =================【產生時間戳記 (sTime)】=================
const sTime = Date.now().toString(); 

// =================【抓取簽章需要的額外欄位】=================
// 從 Header 拿 sVer, sVersion (如果沒有設定，給予空字串，避免拼接時出現 undefined)
const sVer = pm.request.headers.get("ver") || "";
const sVersion = pm.request.headers.get("version") || "";

// 從 URL Params 拿 sFilter, sSorter, sPagination（明文 JSON 字串，跟 body 的 value 一樣要先加密才能上路）
const sFilter = pm.request.url.query.get("filters") || "[]";
const sSorter = pm.request.url.query.get("sorters") || "[]";
const sPagination = pm.request.url.query.get("pagination") || "{}";

// 進行 AES-128-CBC 加密，產生真正會送出的 f, s, p
const encryptedFilter = CryptoJS.AES.encrypt(sFilter, key, {
    iv: iv,
    mode: CryptoJS.mode.CBC,
    padding: CryptoJS.pad.Pkcs7
});
const sF = toBase64Url(encryptedFilter.toString());

const encryptedSorter = CryptoJS.AES.encrypt(sSorter, key, {
    iv: iv,
    mode: CryptoJS.mode.CBC,
    padding: CryptoJS.pad.Pkcs7
});
const sS = toBase64Url(encryptedSorter.toString());

const encryptedPagination = CryptoJS.AES.encrypt(sPagination, key, {
    iv: iv,
    mode: CryptoJS.mode.CBC,
    padding: CryptoJS.pad.Pkcs7
});
const sP = toBase64Url(encryptedPagination.toString());

// 把明文 filter/sorter/pagination 換成加密後的 f/s/p，寫回真正會發送的 URL
pm.request.url.query.remove("filters");
pm.request.url.query.remove("sorters");
pm.request.url.query.remove("pagination");
pm.request.url.query.upsert({ key: "f", value: sF });
pm.request.url.query.upsert({ key: "s", value: sS });
pm.request.url.query.upsert({ key: "p", value: sP });

// =================【計算 Go 原生 MD5 簽章】=================
// 後端拼接順序: sVer, sVersion, sK, sTime, sF, sS, sP, sV, SALT
const aStrings = [sVer, sVersion, sK, sTime, sF, sS, sP, sV, SIGNATURE_SALT];
const sStrings = aStrings.join("|");

// 計算 MD5 雜湊值並轉為小寫 32 位元字串
const sMd5Signature = CryptoJS.MD5(sStrings).toString(CryptoJS.enc.Hex);

// =================【重寫 Body 只發送 {"v": sV}】=================
const oBody = { "v": sV };
pm.request.body.update({
    mode: 'raw',
    raw: JSON.stringify(oBody)
});

// =================【動態寫入 Headers】=================


// 2. 把 sK 塞進 Header K
pm.request.headers.upsert({ key: "K", value: sK });

// 3. 把 時間戳記 塞進 Header time
pm.request.headers.upsert({ key: "time", value: sTime });

// 4. 把 簽章 塞進 Header sign (請根據你後端 GetHeader 的 Key 命名調整，例如 "sign" 或 "signature")
pm.request.headers.upsert({ key: "signature", value: sMd5Signature });

// 1. 強制設定 Content-Type
pm.request.headers.upsert({ key: "Content-Type", value: "application/json" });

// =====================================================

console.log("--- 【混合加密與簽章成功】 ---");
console.log("拼接字串原貌:", sStrings);
console.log("產生的 MD5 簽章 (sign):", sMd5Signature);