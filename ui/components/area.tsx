type AreaProps = Omit<React.ComponentPropsWithoutRef<"div">, "className">;

const Area = (props: AreaProps) => {
  return (
    <div {...props} className="p-4 my-4 border-2 border-pink-600 rounded" />
  );
};

export default Area;
