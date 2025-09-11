import { GiArrowDunk } from "react-icons/gi";
import { Link } from '@tanstack/react-router'

function Logo() {
  return (
    <>
      <Link to="/">
        <div className="flex items-center space-x-2">
          <GiArrowDunk className="w-6 h-6" />
          <p className="font-bold">Tunnel</p>
        </div>
      </Link>
    </>
  )
}
export default Logo;
