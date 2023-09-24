import cls from './home.module.css';

import {useState} from "react";

import {Header} from "../../sections/header/header.tsx";
import {ContentWrapper} from "../../sections/content/content.tsx";

export function Home() {
    const [isSideMenuOpen, setIsSideMenuOpen] = useState(false)
    return (
        <div className={cls.Home}>
            <div className={cls.contentWrap}>
                <ContentWrapper
                    isOpen={isSideMenuOpen}
                />
            </div>

            <Header
                isOpen={isSideMenuOpen}
                setIsOpen={setIsSideMenuOpen}
            />
        </div>
    )
}
