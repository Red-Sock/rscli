import cls from './burger.module.css'
import {motion} from 'framer-motion';

export function Burger(props: { isOpen: boolean; }) {

    return (
        <motion.div
            layout

            className={cls.Burger}
            data-isOpen={props.isOpen}
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
