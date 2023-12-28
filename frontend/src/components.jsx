import styled from "styled-components";

export const StyledApp = styled.div`
  height: 100vh;
  text-align: center;
`;

export const Logo = styled.img`
  display: block;
  width: 20%;
  height: 25%;
  margin: auto;
  padding: 10% 0 0;
  background-position: center;
  background-repeat: no-repeat;
  background-size: 70% 70%;
  background-origin: content-box;
`;

export const Prompt = styled.div`
  border: 1px solid green;
  background-color: lightgreen;
  color: darkgreen;
`;

export const InputBox = styled.div`
  .input {
    border: none;
    border-radius: 3px;
    outline: none;
    height: 30px;
    line-height: 30px;
    padding: 0 10px;
    background-color: rgba(240, 240, 240, 1);
    -webkit-font-smoothing: antialiased;
    &:hover {
      background-color: rgba(255, 255, 255, 1);
    }
    &:focus {
      background-color: rgba(255, 255, 255, 1);
    }
  }
  .btn {
    width: 60px;
    height: 30px;
    line-height: 30px;
    border-radius: 3px;
    border: none;
    margin: 0 0 0 20px;
    padding: 0 8px;
    cursor: pointer;
    &:hover:enabled {
      background-image: linear-gradient(to top, #cfd9df 0%, #e2ebf0 100%);
      color: #333333;
    }
  }
`;

export const Result = styled.div`
  width: 40%;
  height: 20px;
  line-height: 20px;
  margin: auto;
  overflow-wrap: break-word;
`;

export const Header = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  padding: .5rem;
  // border: 1px solid red;
`;

export const HeaderMiddle = styled.div`
  font-size: 1.75em;
  font-weight: bold;
  text-align: center;
  width: 70%;
  // border: 1px solid white;
`;

export const HeaderSide = styled.div`
  text-align: right;
  vertical-align: top;
  width: 15%;
  // border: 1px solid white;
`;

export const Body = styled.div`
  display: flex;
  justify-content: space-between;
  // border: 1px solid red;
`;

export const BodyMiddle = styled.div`
  align-items: flex-start;
  text-align: center;
  width: 75%;
  border: 1px solid white;
`;

export const BodySide = styled.div`
  display: flex;
  vertical-align: center;
  width: 25%;
  border: 1px solid white;
`;

export const Footer = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
  position: fixed;
  bottom: 1rem;
  width: 100vw;
  // border: 1px solid red;
`;

export const FooterMiddle = styled.div`
  font-size: 0.75em;
  text-align: center;
  width: 70%;
  // border: 1px solid white;
`;

export const FooterSide = styled.div`
  text-align: right;
  vertical-align: bottom;
  display: flex;
  align-items: flex-end;
  width: 15%;
  // border: 1px solid white;
`;
