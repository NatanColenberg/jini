import React from "react";

interface Props {
  items: any[];
}

export default function ItemList(props: Props) {
  return (
    <div>
      <h1>Items:</h1>
      {props.items.length ? (
        props.items.map((item) => <h1 key={item.title}>{item.title}</h1>)
      ) : (
        <h3>No Items to Display</h3>
      )}
    </div>
  );
}
