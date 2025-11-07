// Load the Navbar and set the links
document.addEventListener("DOMContentLoaded", function () {
	fetch("../components/navbar.html")
		.then((response) => response.text())
		.then((data) => {
			document.body.insertAdjacentHTML("beforeend", data);

			document.getElementById("nav-home").href = getURL();
			document.getElementById("nav-disease-detection").href = getURL() + "/disease-detection";
			document.getElementById("nav-settings").href = getURL() + "/settings";
		})
		.catch((error) => console.error("Error loading navbar:", error));
});
