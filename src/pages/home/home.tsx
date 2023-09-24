import cls from './home.module.css';

import {useState} from "react";

import {Header} from "../../sections/header/header.tsx";
import {Footer} from "../../sections/footer/footer.tsx";
import {Sidebar} from "../../sections/sidebar/sidebar.tsx";
import {ContentDisplay} from "../../sections/content/content.tsx";

export function Home() {
    const [isSideMenuOpen, setIsSideMenuOpen] = useState(false)
    return (
        <div className={cls.Home}>
            <Header
                isOpen={isSideMenuOpen}
                setIsOpen={setIsSideMenuOpen}
            />

            <div className={cls.contentWrap}>
                <ContentDisplay/>
                <Sidebar
                    isOpen={isSideMenuOpen}
                />
            </div>

            <Footer/>
        </div>
    )
}
