import cls from './header.module.css'

import {Burger} from '../../components/burger/burger.tsx';
import {Sidebar} from "../sidebar/sidebar.tsx";
import {useState} from "react";

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
        </header>
    );
}
