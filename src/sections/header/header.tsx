import cls from './header.module.css'
import {Burger} from "../../components/burger/burger.tsx";

export function Header(props: { setIsOpen: (isOpen: boolean) => void; isOpen: boolean; }) {
    return (
        <header className={cls.Header}>
            <div className={cls.BurgerContainer}
                 onClick={()=> props.setIsOpen(!props.isOpen)}>
                <Burger
                    isOpen={props.isOpen}
                />
            </div>
        </header>
    );
}
