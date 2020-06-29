import React from "react";
import Item from "./item/item";
import "./itemList.css";

interface Props {
  items: any[];
  onItemRemove: (id: string) => void;
}

export default function ItemList(props: Props) {
  return (
    <div className="itemList-container">
      <div className="itemList-title">Items:</div>
      <div className="itemList-list">
        {props.items.length ? (
          props.items.map((item) => <Item id={item.id} title={item.title} onRemove={props.onItemRemove}/>)
        ) : (
          <h3>No Items to Display</h3>
        )}
      </div>
    </div>
  );
}
