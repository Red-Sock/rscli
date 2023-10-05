import cls from './sidebar.module.css'

import { TreeView} from "@mui/x-tree-view";

import {FoldedIcon} from '../../components/tree/foldedIcon';
import {UnfoldedIcon} from '../../components/tree/unfoldedIcon';
import {ThemeProvider, createTheme} from "@mui/material";
import {useState} from "react";
import {Node} from "../../components/tree/node";

export function Sidebar(props: {setContentViaLink(link: string):void}) {

    const THEME = createTheme({
        typography: {
            "fontFamily": `"Comfortaa"`,
            // @ts-ignore
            "fontSize": '2em',
            "fontWeightLight": 300,
            "fontWeightRegular": 400,
            "fontWeightMedium": 500
        }
    });

    const [list, setList] = useState(
        [
            {
                name: "home",
                childs:
                    [
                        {
                            name: "uber style guid",
                            childs: [
                                {
                                    name: "EN",
                                    link: "https://raw.githubusercontent.com/uber-go/guide/master/style.md"
                                },
                                {
                                    name: "RU",
                                    link: "https://raw.githubusercontent.com/sau00/uber-go-guide-ru/master/style.md"
                                }
                            ]
                        }
                    ]
            }
        ]
    )

    let nodeId = 1;
    function getNextNodeId() :string {
        return (nodeId++).toString()
    }

    return (
        <div className={cls.Sidebar}>
            <div className={cls.Content}>
                <ThemeProvider theme={THEME}>
                    <TreeView
                        aria-label="Documentation"
                        defaultCollapseIcon={<FoldedIcon/>}
                        defaultExpandIcon={<UnfoldedIcon/>}
                    >
                        {list.map((elem: any)=>{
                            return <Node
                                nextNodeId = {getNextNodeId}
                                openLink={props.setContentViaLink}
                                node={elem}
                            />})}
                    </TreeView>
                </ThemeProvider>
            </div>
        </div>
    );
}
