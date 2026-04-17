export async function getDeviceData() {
    const nav = navigator

    return {
        userAgent: nav.userAgent,
        platform: nav.platform,
        language: nav.language
    }
}