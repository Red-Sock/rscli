import { motion } from 'framer-motion';
import cls from './sidebar.module.css'
import {TreeItem, TreeView} from "@mui/x-tree-view";

import {FoldedIcon} from '../../components/tree/foldedIcon';
import {UnfoldedIcon} from '../../components/tree/unfoldedIcon';
export function Sidebar (props: { isOpen: any; }) {

    return (
        <motion.div
            className={cls.Sidebar}
            layout
            data-isOpen={props.isOpen}
        >
            <div className={cls.Content}>
               <TreeView
                   aria-label="Documentation"
                   defaultCollapseIcon={<FoldedIcon />}
                   defaultExpandIcon={<UnfoldedIcon />}
                   sx={{ height: 240, flexGrow: 1, maxWidth: 400, overflowY: 'auto' }}
               >
                   <TreeItem nodeId="1" label="Applications">
                       <TreeItem nodeId="2" label="Calendar" />
                   </TreeItem>
                   <TreeItem nodeId="5" label="Documents">
                       <TreeItem nodeId="10" label="OSS" />
                       <TreeItem nodeId="6" label="MUI">
                           <TreeItem nodeId="8" label="index.js" />
                       </TreeItem>
                   </TreeItem>
               </TreeView>
            </div>
        </motion.div>
    );
}
