const path = "https://raw.githubusercontent.com/Red-Sock/rscli/docs/docs";
const idx = window.location.pathname.indexOf("/", 2)

export function getResourceURLs(url: string, setContent: (value: string) => void, setNewUrl: (arg0: string) => void) {
    fetch(url).then(async response => {
        if (response.ok) {
            setContent(await response.text())
        } else {
            setNewUrl(window.location.protocol+"//"+window.location.host+"/rscli/home")
        }
    })
}

function clearMd(req: string) : string {
    let idx = req.indexOf("<!-- ")
    let iterations = 10
    while (idx !== -1 && iterations > 0) {
        iterations--
        let idxEnd = req.indexOf("-->")
        if (idxEnd == -1) {
            idxEnd=idx+4
        }

        req = req.substring(idx, idxEnd)
    }
    return req
}

export function getGithubResourceURLs( setContent: (value: string) => void, setNewUrl: (arg0: string) => void) {
    fetch(path + window.location.pathname.substring(idx)+".md").then(async response => {
        if (response.ok) {
            setContent(await response.text())
        } else {
            setNewUrl(window.location.protocol+"//"+window.location.host+"/rscli/home")
        }
    })
}
