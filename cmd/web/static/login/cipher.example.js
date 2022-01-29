// generates a 128 bit key for encryption and decryption in base64 format
const generateKey = async () => {
    const key = await window.crypto.subtle.generateKey({
            name: 'AES-GCM',
            length: 128
        },
        true, [
            'encrypt',
            'decrypt'
        ]);

    const exportedKey = await window.crypto.subtle.exportKey(
        'raw',
        key,
    );
    return bufferToBase64(exportedKey);
}

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

const cryptGcm = async (base64Key, bytes) => {
    const key = await loadKey(base64Key);
    const iv = window.crypto.getRandomValues(new Uint8Array(12));
    const algorithm = {
        iv,
        name: 'AES-GCM'
    };
    const cipherData = await window.crypto.subtle.encrypt(
        algorithm,
        key,
        bytes
    );

    // prepend the random IV bytes to raw cipherdata
    const cipherText = concatArrayBuffers(iv.buffer, cipherData);
    return bufferToBase64(cipherText);
}

// concatenate two array buffers
const concatArrayBuffers = (buffer1, buffer2) => {
    let tmp = new Uint8Array(buffer1.byteLength + buffer2.byteLength);
    tmp.set(new Uint8Array(buffer1), 0);
    tmp.set(new Uint8Array(buffer2), buffer1.byteLength);
    return tmp.buffer;
}

let plaintext = "couple robot escape silent main once smoke check good basket mimic similar";
let plaintextBytes = (new TextEncoder()).encode(plaintext, 'utf-8');
let encryptionKey = "FjTOwFMyBNFpHS4YGXfitm7+V31BrYVEUNQXuOnlcpI=";
let ciphertext = await cryptGcm(encryptionKey, plaintextBytes);


let plaintext = "123456";
let plaintextBytes = (new TextEncoder()).encode(plaintext, 'utf-8');
console.log(plaintextBytes);
let encryptionKey = "FjTOwFMyBNFpHS4YGXfitm7+V31BrYVEUNQXuOnlcpI=";
let ciphertext = cryptGcm(encryptionKey, plaintextBytes);
console.log("plaintext: ", plaintext);
console.log("encryptionKey (base64):", encryptionKey);
console.log("ciphertext (base64):", ciphertext);