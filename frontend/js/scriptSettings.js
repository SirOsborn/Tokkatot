// DOM Elements
const autoModeToggle = document.getElementById("autoModeToggle");
// const scheduleModeToggle = document.getElementById("scheduleModeToggle");
const conveyerToggle = document.getElementById("conveyerToggle");
const fanToggle = document.getElementById("fanToggle");
const lightToggle = document.getElementById("lightToggle");
const feederToggle = document.getElementById("feederToggle");
const pumpToggle = document.getElementById("pumpToggle");
// const feedingTimesContainer = document.getElementById("feedingTimes");
// const saveScheduleButton = document.getElementById("saveSchedule");
const notification = document.getElementById("notification");

// Initialize the system on page load
document.addEventListener("DOMContentLoaded", () => {
    // Fetch and display initial settings
    fetchInitialSettings();

    // Attach event listeners for toggles
    autoModeToggle.addEventListener("change", async () => {
        const state = autoModeToggle.checked;
        await handleModeToggle("/api/toggle-auto", state);
    });

    //scheduleModeToggle.addEventListener("change", () => {
    /*    const state = scheduleModeToggle.checked;
        handleModeToggle("/api/toggle-schedule", state);
        if (state) {
            showNotification("Schedule Mode enabled. Configure your settings.", "info");
        } else {
            showNotification("Schedule Mode disabled.", "info");
        }
    }); */

    // Attach event listeners for immediate toggles
    conveyerToggle.addEventListener("change", () =>
        handleImmediateToggle("/api/toggle-belt", conveyerToggle.checked),
    );
    fanToggle.addEventListener("change", () =>
        handleImmediateToggle("/api/toggle-fan", fanToggle.checked),
    );
    lightToggle.addEventListener("change", () =>
        handleImmediateToggle("/api/toggle-bulb", lightToggle.checked),
    );
    feederToggle.addEventListener("change", () =>
        handleImmediateToggle("/api/toggle-feeder", feederToggle.checked),
    );
    pumpToggle.addEventListener("change", () =>
        handleImmediateToggle("/api/toggle-pump", pumpToggle.checked),
    );

    // Attach event listener for saving schedule settings
  /*  saveScheduleButton.addEventListener("click", saveScheduleSettings);

    // Launch schedule modal (if existing modal functionality exists)
    document
        .getElementById("addFeedingTime")
        .addEventListener("click", () => addFeedingTimeInput("")); */
});

// Fetch initial settings state from the backend
async function fetchInitialSettings() {
    try {
        const response = await fetch("/api/get-initial-state");
        if (!response.ok) {
            throw new Error("Failed to fetch initial state.");
        }

        let data = await response.json();
        data = JSON.parse(data.data);
        // console.log("Fetched initial state:", data);

        updateUI(data);
    } catch (error) {
        console.error("Error fetching initial state:", error);
        showNotification("Unable to load initial settings!", "error");
    }
}

setInterval(fetchInitialSettings, 100);

// Update the UI based on retrieved settings
function updateUI(data) {
    // Update mode toggles
    autoModeToggle.checked = data.auto_mode;
    // scheduleModeToggle.checked = data.scheduleMode;

    // Update immediate toggles
    conveyerToggle.checked = data.conveyer;
    fanToggle.checked = data.fan;
    lightToggle.checked = data.bulb;
    feederToggle.checked = data.feeder;
    pumpToggle.checked = data.pump;

    // Update schedule settings
  /*  document.getElementById("lightStart").value =
        data.schedule.lighting.start || "06:00";
    document.getElementById("lightEnd").value =
        data.schedule.lighting.end || "18:00";

    // Update feeding times schedule
    feedingTimesContainer.innerHTML = ""; // Clear existing feeding times
    data.schedule.feeding.forEach((time) => addFeedingTimeInput(time));

    // Update water interval
    document.getElementById("waterInterval").value =
        data.schedule.waterInterval || 60;

    // Update environmental thresholds
    document.getElementById("tempMin").value =
        data.schedule.tempThreshold.min || 20;
    document.getElementById("tempMax").value =
        data.schedule.tempThreshold.max || 25;
    document.getElementById("humidityMin").value =
        data.schedule.humThreshold.min || 40;
    document.getElementById("humidityMax").value =
        data.schedule.humThreshold.max || 60; */
}

// Handle toggling Auto Mode
async function handleModeToggle(endpoint, state) {
    try {
        const response = await fetch(endpoint, {
            method: "GET",
        });

        if (!response.ok) {
            throw new Error(`Failed to toggle ${endpoint}.`);
        }

        let result = await response.json();
        result = JSON.parse(result.state);
        console.log(`Toggled ${endpoint}: `, result);

        if (!result) {
            conveyerToggle.checked = false;
            fanToggle.checked = false;
            lightToggle.checked = false;
            feederToggle.checked = false;
            pumpToggle.checked = false;
        }
    } catch (error) {
        console.error(`Error toggling ${endpoint}:`, error);
        showNotification("Failed to update Auto Mode toggle.", "error");

        // Revert the toggle state on error
        autoModeToggle.checked = !state;
    }
}

// Handle immediate toggles (like belt, fan, light, feeder, and water)
async function handleImmediateToggle(endpoint, state) {
    try {
        const response = await fetch(endpoint, {
            method: "GET",
        });

        if (!response.ok) {
            throw new Error(`Failed to toggle ${endpoint}.`);
        }

        let result = await response.json();
        result = JSON.parse(result.state);
        console.log(`Toggled ${endpoint}: `, result);
    } catch (error) {
        console.error(`Error toggling ${endpoint}:`, error);
        showNotification("Failed to update device state.", "error");

        // Revert the toggle state on error
        const device = endpoint.split("-")[1];
        document.getElementById(`${device}Toggle`).checked = !state;
    }
}

// Handle saving updated schedule settings
/* async function saveScheduleSettings() {
    try {
        // Gather lighting schedule inputs
        const lightStart = document.getElementById("lightStart").value;
        const lightEnd = document.getElementById("lightEnd").value;

        // Gather feeding schedule inputs
        const feedingTimes = Array.from(
            feedingTimesContainer.querySelectorAll("input[type='time']"),
        ).map((input) => input.value);

        // Gather water interval input
        const waterInterval = parseInt(
            document.getElementById("waterInterval").value,
            10,
        );

        // Gather environmental thresholds
        const tempMin = parseFloat(document.getElementById("tempMin").value);
        const tempMax = parseFloat(document.getElementById("tempMax").value);
        const humidityMin = parseFloat(
            document.getElementById("humidityMin").value,
        );
        const humidityMax = parseFloat(
            document.getElementById("humidityMax").value,
        );

        // Create payload
        const payload = {
            lighting: {
                start: lightStart,
                end: lightEnd,
            },
            feeding: feedingTimes,
            waterInterval: waterInterval,
            tempThreshold: {
                min: tempMin,
                max: tempMax,
            },
            humThreshold: {
                min: humidityMin,
                max: humidityMax,
            },
        };

        console.log("Saving schedule settings with payload:", payload);

        const response = await fetch("/api/schedule", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(payload),
        });

        if (!response.ok) {
            throw new Error("Failed to save the schedule settings.");
        }

        const result = await response.json();
        console.log("Schedule saved successfully:", result);

        showNotification("Schedule updated successfully!", "success");
    } catch (error) {
        console.error("Error saving schedule:", error);
        showNotification("Failed to save schedule settings.", "error");
    }
}

// Helper function to dynamically add a feeding time row
 function addFeedingTimeInput(time = "") {
    const feedingTimeDiv = document.createElement("div");
    feedingTimeDiv.classList.add("feeding-time");

    const input = document.createElement("input");
    input.type = "time";
    input.value = time;

    const removeButton = document.createElement("button");
    removeButton.textContent = "លុប"; // "Remove" in Khmer
    removeButton.classList.add("remove-time");
    removeButton.addEventListener("click", () => feedingTimeDiv.remove());

    feedingTimeDiv.appendChild(input);
    feedingTimeDiv.appendChild(removeButton);
    feedingTimesContainer.appendChild(feedingTimeDiv);
} */

// Helper function to display notifications
function showNotification(message, type) {
    notification.textContent = message;
    notification.className = `notification ${type}`;
    setTimeout(() => {
        notification.className = "notification"; // Clear notification
    }, 3000);
}

// Add periodic polling to sync the state of the toggles
// setInterval(async () => {
//     try {
//         const response = await fetch("/api/get-current-data");
//         if (!response.ok) {
//             throw new Error("Failed to fetch current data.");
//         }

//         let data = await response.json();
//         data = JSON.parse(data.data);
//         console.log("Fetched current data:", data);

//         // Update the UI based on the fetched data
//         updateUI(data);
//     } catch (error) {
//         console.error("Error fetching current data:", error);
//     }
// }, 3000); // Poll every 3 seconds