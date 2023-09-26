import cls from './home.module.css';

import {Header} from "../../sections/header/header";
import {ContentWrapper} from "../../sections/content/content";

export function Home() {
    return (
        <div className={cls.Home}>

            <div className={cls.headerWrap}>
                <Header/>
            </div>

            <div className={cls.contentWrap}>
                <ContentWrapper/>
            </div>


        </div>
    )
}
