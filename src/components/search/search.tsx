import cls from "./search.module.css"

import {memo} from "react";

export const Search = memo(() => {
        return (
            <div className={cls.SearchBox}>
                <link href="https://fonts.googleapis.com/icon?family=Material+Icons" rel="stylesheet"/>
                <button className={cls.btnSearch}>
                    <i className="material-icons">search</i>
                </button>
                <input type="text" className={cls.inputSearch} placeholder="Type to Search..."/>
            </div>
        )
    }
)
