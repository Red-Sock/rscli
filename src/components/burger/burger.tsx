import cls from './burger.module.css'
import {motion} from "framer-motion";

export function Burger(props: { isOpen: boolean; }) {
    return (
        <motion.div
            className={cls.Burger}
            data-isOpen={props.isOpen}
            layout
        >
            <motion.div
                data-isOpen={props.isOpen}
                className={cls.TopLine}
            />
            <motion.div
                data-isOpen={props.isOpen}
                className={cls.MidLine}/>
            <motion.div
                data-isOpen={props.isOpen}
                className={cls.BotLine}/>
        </motion.div>
    );
}
