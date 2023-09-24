import cls from './header.module.css'

import { motion } from 'framer-motion';

import {Burger} from '../../components/burger/burger.tsx';

export function Header(props: { setIsOpen: (isOpen: boolean) => void; isOpen: boolean; }) {
    return (
        <header className={cls.Header}>
            <motion.div
                layout
                className={cls.BurgerContainer}
                data-isOpen={props.isOpen}
                 onClick={()=> props.setIsOpen(!props.isOpen)}>
                <Burger
                    isOpen={props.isOpen}
                />
            </motion.div>
        </header>
    );
}
