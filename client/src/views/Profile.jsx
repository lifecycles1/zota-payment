import { useParams } from "react-router-dom";

const Profile = () => {
  const { id } = useParams();

  return (
    <div>
      <h1>User Profile</h1>
      <p>UID: {id}</p>
      {/* Display more user account or deposit information */}
    </div>
  );
};

export default Profile;
