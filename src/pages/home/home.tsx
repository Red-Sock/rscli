import cls from './home.module.css';

import {Header} from "../../sections/header/header";
import {ContentWrapper} from "../../sections/content/content";
import {Footer} from "../../sections/footer/footer";
import {Sidebar} from "../../sections/sidebar/sidebar";

export function Home() {
    return (
        <div className={cls.Home}>

            <div className={cls.headerWrap}>
                <Header/>
            </div>

            <div className={cls.contentWrap}>
                <ContentWrapper/>
            </div>

            <div className={cls.sideMenuWrap}>
                <Sidebar/>
            </div>

        </div>
    )
}
