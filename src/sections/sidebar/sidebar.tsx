import cls from './sidebar.module.css'

import {TreeItem, TreeView} from "@mui/x-tree-view";

import {FoldedIcon} from '../../components/tree/foldedIcon';
import {UnfoldedIcon} from '../../components/tree/unfoldedIcon';

export function Sidebar() {

    return (
        <div
            className={cls.Sidebar}
        >
            <div className={cls.Content}>
                <TreeView
                    aria-label="Documentation"
                    defaultCollapseIcon={<FoldedIcon/>}
                    defaultExpandIcon={<UnfoldedIcon/>}
                    sx={{
                        height: '100%',
                        flexGrow: 1,
                        overflowY: 'auto',
                        "& .MuiTreeItem-label": {
                            fontSize: '1em',
                        }
                    }}
                >
                    <TreeItem nodeId="1" label="Applications">
                        <TreeItem nodeId="2" label="Calendar"/>
                    </TreeItem>

                    <TreeItem nodeId="3" label="Documents">
                        <TreeItem nodeId="4" label="OSS"/>
                        <TreeItem nodeId="5" label="MUI">
                            <TreeItem nodeId="6" label="index.js"/>
                        </TreeItem>
                    </TreeItem>

                    <TreeItem nodeId="7" label="Applications">
                        <TreeItem nodeId="8" label="Calendar"/>
                    </TreeItem>

                    <TreeItem nodeId="9" label="Documents">
                        <TreeItem nodeId="10" label="OSS"/>
                        <TreeItem nodeId="11" label="MUI">
                            <TreeItem nodeId="12" label="index.js"/>
                        </TreeItem>
                    </TreeItem>

                    <TreeItem nodeId="13" label="Applications">
                        <TreeItem nodeId="14" label="Calendar"/>
                    </TreeItem>

                    <TreeItem nodeId="15" label="Documents">
                        <TreeItem nodeId="16" label="OSS"/>
                        <TreeItem nodeId="17" label="MUI">
                            <TreeItem nodeId="18" label="index.js"/>
                        </TreeItem>
                    </TreeItem>
                </TreeView>
            </div>
        </div>
    );
}
