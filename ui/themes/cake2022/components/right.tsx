interface RightProps {
  children: React.ReactChild;
}

const Right = ({ children }: RightProps) => {
  return (
    <div
      style={{ display: "flex", flexDirection: "row-reverse", width: "100%" }}
    >
      {children}
    </div>
  );
};
export default Right;
