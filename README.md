# Distributed Database System (Go + MariaDB + Streamlit)

## 📌 Project Overview

This project implements a **Distributed Database System** using **Go**, **MariaDB**, **mysql**, and **Python (Streamlit)**. The system follows a **Master-Snap architecture** where:

- The **Master node** acts as the central authority:
  - Handles database creation, table creation, and deletion.
  - Manages centralized query execution.
  - Synchronizes data with Snap nodes.
  - Hosts a **Streamlit GUI** to send SQL commands and view results.

- The **Snap nodes** are terminal-based clients:
  - Interact directly with users via a CLI.
  - Support **INSERT**, **UPDATE**, **DELETE**, **SELECT ALL**, and **SEARCH** operations.
  - Perform **bidirectional replication**, meaning any change on a Snap node is:
    - Executed locally,
    - Sent to the Master,

---

## ⚙️ Features

### ✅ Master Node
- Creates the `MOST` database and tables.
- Deletes databases/tables (only allowed on Master).
- Executes and logs queries.
- Hosts a **Streamlit interface** to allow SQL command input.
- Can **export** the entire database as a `.sql` file using `mysqldump`.
- Sends `.sql` snapshots to Snap nodes upon request (initial sync).

### ✅ Snap Node
- Terminal-based interface with the following operations:
  - `1 - INSERT`
  - `2 - UPDATE`
  - `3 - DELETE`
  - `4 - SELECT ALL`
  - `5 - SEARCH`
  - `6 - Exit`
- Allows dynamic table selection with a **Back** option to return to the main menu.
- Automatically receives the database dump from the Master on first connection.
- Sends every DML operation to Master and other Snaps for **bidirectional sync**.

---

## 🔄 Replication Model

- **Initial Sync**: On startup, a Snap node connects to the Master and receives a `.sql` dump, which is imported locally.
- **Bidirectional Replication**:
  - Any INSERT, UPDATE, or DELETE operation on **any node** (Master or Snap):
    - Executes locally.
    - Is sent to the Master.
    - Is broadcasted to all connected Snap nodes.

---

## 🛠️ Technologies Used
- **Go**: Core language for Master and Snap communication.
- **MySQL**: Database management system for data storage.
- **Streamlit (Python)**: For GUI interface on Master.
- **TCP Sockets**: Used for communication between Master and Snap nodes.

---

## 🗂️ Project Structure




                        +--------------------------+
                        |      Master Node         |
                        |   (Linux + Streamlit GUI) |
                        |--------------------------|
                        | - MariaDB Database       |
                        | - Query Execution        |
                        | - Database Dump Export   |
                        | - Sync with Snap Nodes   |
                        +--------------------------+
                                |
                      +----------------------------+
                      | TCP Communication           |
                      +----------------------------+
                                |
    +---------------------------+---------------------------+
    |                           |                           |
+-----------+             +-----------+             +-----------+
|  Snap 1   |             |  Snap 2   |             |  Snap 3   |
|  (Terminal)|             |  (Terminal)|             |  (Terminal)|
+-----------+             +-----------+             +-----------+
    |                           |                           |
+-----------------+        +-----------------+        +-----------------+
| - INSERT        |        | - INSERT        |        | - INSERT        |
| - UPDATE        |        | - UPDATE        |        | - UPDATE        |
| - DELETE        |        | - DELETE        |        | - DELETE        |
| - SELECT ALL    |        | - SELECT ALL    |        | - SELECT ALL    |
| - SEARCH        |        | - SEARCH        |        | - SEARCH        |
+-----------------+        +-----------------+        +-----------------+
    |                           |                           |
    +---------------------------+---------------------------+
                                |
                          +--------------------------+
                          | Bidirectional Sync       |
                          | (Master <-> Snap Nodes)  |
                          +--------------------------+
                                              
                             
