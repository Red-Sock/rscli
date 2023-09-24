import { motion } from 'framer-motion';
import cls from './sidebar.module.css'
export function Sidebar (props: { isOpen: any; }) {
    return (
        <motion.div
            className={cls.Sidebar}
            layout
            data-isOpen={props.isOpen}
        />
    );
}
