import Sidebar from "../../components/MapEditor/Sidebar";
import Map from "../../components/MapEditor/Map";

function Homepage() {
  return (
    <div className="container-fluid px-0" style={{ overflow: "hidden" }}>
      <div className="row no-gutters" style={{ height: "calc(100vh - 64px)" }}>
        {/* Sidebar, tall, sticky, nav-tabs, content */}
        <div className="col-3 bg-light border-end p-0">
          <Sidebar />
        </div>
        {/* Map */}
        <div className="col-9 p-0">
          <div style={{ height: "100%", width: "100%" }}>
            <Map />
          </div>
        </div>
      </div>
    </div>
  );
}

export default Homepage;