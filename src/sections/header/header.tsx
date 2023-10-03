import cls from './header.module.css'

import {Search} from "../../components/search/search";

export function Header() {
    return (
        <header className={cls.Header}>
            <div className={cls.SearchContainer}>
                <Search/>
            </div>
        </header>
    );
}
