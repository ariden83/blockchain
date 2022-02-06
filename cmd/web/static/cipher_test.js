// arrayBuffer to base64
const bufferToBase64 = (arrayBuffer) => {
    return window.btoa(String.fromCharCode(...new Uint8Array(arrayBuffer)));
}

// load a base64 encoded key
const loadKey = async (base64Key) => {
    return await window.crypto.subtle.importKey(
        'raw',
        base64ToBuffer(base64Key),
        "AES-GCM",
        true, [
            "encrypt",
            "decrypt"
        ]
    );
}

// base64 to arrayBuffer
const base64ToBuffer = (base64) => {
    const binary_string = window.atob(base64);
    const len = binary_string.length;
    let bytes = new Uint8Array(len);
    for (let i = 0; i < len; i++) {
        bytes[i] = binary_string.charCodeAt(i);
    }
    return bytes.buffer;
}

const cryptGcm = async (base64Key, seed) => {
    const bytes = (new TextEncoder()).encode(seed, 'utf-8');
    const key = await loadKey(base64Key);
    const iv = window.crypto.getRandomValues(new Uint8Array(12));
    const cipherData = await window.crypto.subtle.encrypt(
        {
            iv,
            name: 'AES-GCM'
        },
        key,
        bytes
    );
    const cipherText = concatArrayBuffers(iv.buffer, cipherData);
    return bufferToBase64(cipherText);
}

const DecryptGcm = async (base64Key, cipherTextBase64) => {
    const cipherText = base64ToBuffer(cipherTextBase64)
    const key = await loadKey(base64Key);
    const data = ArrayBuffersDecoder(cipherText);

    const decrypted = await window.crypto.subtle.decrypt(
        {
            iv: data.iv,
            name: 'AES-GCM'
        },
        key,
        data.cipher,
    );

    const decoder = new TextDecoder();
    const plaintext = decoder.decode(decrypted);
    return plaintext;
}

// concatenate two array buffers
const concatArrayBuffers = (buffer1, buffer2) => {
    let tmp = new Uint8Array(buffer1.byteLength + buffer2.byteLength);
    tmp.set(new Uint8Array(buffer1), 0);
    tmp.set(new Uint8Array(buffer2), buffer1.byteLength);
    return tmp.buffer;
}

const ArrayBuffersDecoder = (buffer) => {
    let iv = new Uint8Array(buffer.slice(0, 12));
    let cipher = new Uint8Array(buffer.slice(12, buffer.length));
    return {
        iv: iv,
        cipher: cipher,
    }
}

const encrypt = async (base64Key, code) => {
    let e = await cryptGcm(base64Key, code);
    return e
}

const decrypt = async (base64Key, seed) => {
    let e = await DecryptGcm(base64Key, seed);
    return e;
}

const cipher_coder_decoder = async (base64Key, seed) => {
    return encrypt(base64Key, seed).then(cipherTextBase64 => decrypt(base64Key, cipherTextBase64));
}

const cipher_coder_decoder_decompress = async (base64Key, seed) => {
    /************************ encode ********************************/
    const bytes = (new TextEncoder()).encode(seed, 'utf-8');
    let key = await loadKey(base64Key);
    const iv = window.crypto.getRandomValues(new Uint8Array(12));
    const cipherData = await window.crypto.subtle.encrypt(
        {
            iv,
            name: 'AES-GCM'
        },
        key,
        bytes
    );
    let cipherText = concatArrayBuffers(iv.buffer, cipherData);
    let cipherTextBase64 = bufferToBase64(cipherText);
    /************************ decode ********************************/
    return DecryptGcm(base64Key, cipherTextBase64)
}

let data_to_test = [
    "1234556",
    "couple robot escape silent main once smoke check good basket mimic similar"
];

const test_cipher_coder_decoder_decompress = async () => {
    let base64Key = "MYKS4F9T28bF1WJy/FxeGJ7JfkeTSBPhK5vOy/mZuMw=";
    try {
        for (const data of data_to_test) {
            if (data === cipher_coder_decoder(base64Key, data)) {
                throw Error('Assert failed: ' + (data || ''));
            } else {
                console.log("assert ok: "+ (data || ''));
            }
            if (data === cipher_coder_decoder_decompress(base64Key, data)) {
                throw Error('Assert failed: ' + (data || ''));
            } else {
                console.log("assert ok: "+ (data || ''));
            }
        }
    } catch (e) {
        console.log("error :: ", e);
    }
}

test_cipher_coder_decoder_decompress();
