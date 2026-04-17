import "htmx.org"
import { getOrCreateDeviceId, getBrowserInfo } from "./device_data.js"

const deviceId = getOrCreateDeviceId();
const browserInfo = getBrowserInfo();

console.log("device_id:", deviceId);
console.log("browser info:", browserInfo);
