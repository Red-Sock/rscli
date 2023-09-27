const contentSourceURL = "https://raw.githubusercontent.com/Red-Sock/rscli/docs/docs/";
const projectRoot = "/rscli/";
const resourceDelimiter = "#/"
const defaultURI = "home";
const fileExtension = ".md";

export function doStaff() {
    let currentPage = window.location.href
    let correctPageBase = window.location.origin + projectRoot + resourceDelimiter;


    const out = {
        refreshURl: correctPageBase + defaultURI,
        resourceURL: "",
    }


    if (!currentPage.startsWith(correctPageBase)) {
        return out
    }

    let pathSeparatorIdx = currentPage.indexOf(resourceDelimiter);
    if (pathSeparatorIdx === -1) {
        return out
    }

    const resourceURI = currentPage.substring(pathSeparatorIdx + resourceDelimiter.length)
    out.resourceURL = contentSourceURL + "/" + resourceURI + fileExtension;

    return out
}
