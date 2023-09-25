import 'tree.module.css'

import React from "react";


export function Tree(props: {
    data: any;
    toggled: boolean;
    name: string | null;

    isChildElement: boolean;
    isParentToggled: boolean;
    isLast: any;

}) {
    const isDataArray = Array.isArray(props.data);
    const [isToggled, setIsToggled] = React.useState(props.toggled);
    return (
        <div
            className={`tree-element ${props.isParentToggled && 'collapsed'} ${
                props.isChildElement && 'is-child'
            }`}
        >
      <span
          className={isToggled ? 'toggler' : 'toggler closed'}
          onClick={() => setIsToggled(!isToggled)}
      />
            {props.name ? <strong>&nbsp;&nbsp;{props.name}: </strong> : <span>&nbsp;&nbsp;</span>}
            {isDataArray ? '[' : '{'}
            {!isToggled && '...'}
            {Object.keys(props.data).map((v, i, a) =>
                typeof props.data[v] === 'object' ? (
                    <Tree
                        key={`${props.name}-${v}-${i}`}
                        data={props.data[v]}
                        isLast={i === a.length - 1}
                        name={isDataArray ? null : v}
                        toggled={false}
                        isChildElement
                        isParentToggled={props.isParentToggled && isToggled}
                    />
                ) : (
                    <p
                        key={`${props.name}-${v}-${i}`}
                        className={isToggled ? 'tree-element' : 'tree-element collapsed'}
                    >
                        {isDataArray ? '' : <strong>{v}: </strong>}
                        {props.data[v]}
                        {i === a.length - 1 ? '' : ','}
                    </p>
                )
            )}
            {isDataArray ? ']' : '}'}
            {!props.isLast ? ',' : ''}
        </div>
    )
}

