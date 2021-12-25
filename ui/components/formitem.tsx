import React from "react";

interface FormItemProps {
  children: React.ChildNode;
}

const FormItem = ({ children }: FormItemProps) => {
  return <div className="mb-4">{children}</div>;
};

export default FormItem;
