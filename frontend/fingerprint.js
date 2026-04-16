import FingerprintJS from "@fingerprintjs/fingerprintjs"

export async function getVisitorId() {
    const fp = await FingerprintJS.load()
    const result = await fp.get()
    return result.visitorId
}

; (async () => {
    const visitorId = await getVisitorId()

    const response = await fetch("/api/check-vote", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ visitorId })
    })

    const { isVotedToday } = await response.json()
    console.log(isVotedToday)
})()
