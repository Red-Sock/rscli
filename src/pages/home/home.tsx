import cls from './home.module.css';

import {Header} from "../../sections/header/header";
import {ContentWrapper} from "../../sections/content/content";

import {Sidebar} from "../../sections/sidebar/sidebar";
import {Route, Routes} from "react-router-dom";
import {useState} from "react";
import {getResourceURLs} from "../../services/file-fetcher";

export function Home() {

    const [pageContent, setPageContent] = useState(`
# Click menu on the left to select topic ->
`)

    function getContentViaLink(link: string) {
        getResourceURLs(link, setPageContent, (url: string)=> { window.location.replace(url)})
    }

    return (
        <>
            <div className={cls.headerWrap}>
                <Header/>
            </div>

            <div className={cls.Home}>

                <div className={cls.contentWrap}>
                    <Routes>
                        <Route path={"/*"} element={<ContentWrapper content={pageContent}/>}/>
                    </Routes>
                </div>

                <div className={cls.sideMenuWrap}>
                    <Sidebar setContentViaLink={getContentViaLink}/>
                </div>

            </div>
        </>
    )
}
