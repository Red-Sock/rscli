import {TreeItem} from "@mui/x-tree-view";

export function Node(props: {
                         nextNodeId: () => string,
                         openLink: (link: string) => void,
                         node: { name: string, link: string, childs: any[] }
                     }
) {

    return (
        <>
            <TreeItem nodeId={props.nextNodeId()}
                      label={props.node.name}
                      onClick={() => {
                          if (props.node.link.length!=0){
                            props.openLink(props.node.link)
                          }
                      }
                      }>
                {
                    props.node.childs &&
                    props.node.childs.length > 0 ?
                        props.node.childs.map((c: any) => {
                            return <Node nextNodeId={props.nextNodeId}
                                         openLink={props.openLink}
                                         node={c}

                            />
                        }) : ""
                }
            </TreeItem>
        </>
    )
}

