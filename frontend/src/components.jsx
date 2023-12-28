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
  width: 40%;
  line-height: 20px;
  margin: 1.5rem auto;
  nowrap: true;
`;

export const InputBox = styled.div`
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
`;

export const Result = styled.div`
  width: 40%;
  height: 20px;
  line-height: 20px;
  margin: 1.5rem auto;
  overflow-wrap: break-word;
`;

export const Header = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  padding: 1rem;
`;

export const HeaderSection = styled.div`
  flex-basis: ${(props) => props.width};
`;

export const HeaderMiddle = styled(HeaderSection).attrs({ width: "50%" })`
  font-size: 1.75em;
  font-weight: bold;
  text-align: center;
  // border: 1px solid white;
`;

export const HeaderSide = styled(HeaderSection).attrs({ width: "25%" })`
  text-align: right;
  vertical-align: top;
  height: 100%;
  // border: 1px solid white;
`;
