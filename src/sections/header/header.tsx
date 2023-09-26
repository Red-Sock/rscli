import cls from './header.module.css'

import {Burger} from '../../components/burger/burger';
import {Sidebar} from "../sidebar/sidebar";
import {useState} from "react";
import {Search} from "../../components/search/search";

export function Header() {
    const [isSideMenuOpen, setIsSideMenuOpen] = useState(false)

    return (
        <header className={cls.Header}>
            <div
                className={cls.BurgerContainer}
                 onClick={()=> setIsSideMenuOpen(!isSideMenuOpen)}>
                <Burger
                    isOpen={isSideMenuOpen}
                />
            </div>
            <Sidebar
                isOpen={isSideMenuOpen}
            />
            <div className={cls.SearchContainer}>
                <Search/>
            </div>
        </header>
    );
}
