# Booklet: Project Estimation

## 1. Introduction

This document outlines the estimation of software development complexity, effort, and time for the GeoPathPlanner project. We use two standard methods: Function Point Analysis (FPA) to determine the functional size of the application, and the COCOMO II model to estimate the person-months and development time required.

## 2. Function Point Analysis (FPA)

Function Point Analysis measures the size of the software by quantifying the functionality provided to the user. We identify five types of functional components:
*   **External Inputs (EI):** User inputs that add or change data.
*   **External Outputs (EO):** System outputs to the user (reports, screens).
*   **External Inquiries (EQ):** User requests that retrieve data without changing it.
*   **Internal Logical Files (ILF):** Data stored and maintained within the system.
*   **External Interface Files (EIF):** Data referenced by the application but maintained by another system.

### 2.1. Component Identification and Complexity

First, we map our user stories to these functional components and classify their complexity (Low, Average, High).

| User Story ID(s) | Component Description | Type | Complexity |
| :--------------- | :-------------------- | :--- | :--------- |
| 7, 8, 9, 10      | User Registration/Login/Profile | EI   | Average    |
| 2, 3, 18         | Define Route Inputs (Map/File)  | EI   | High       |
| 16               | View Route History              | EQ   | Average    |
| 5                | Display Computed Route          | EO   | Average    |
| 11               | Show Routing Error Message      | EO   | Low        |
| 6                | Download/Export Route           | EO   | Low        |
| -                | User Account Data               | ILF  | Average    |
| -                | Route History Data              | ILF  | Average    |
| -                | Geo-data Service (Nominatim)    | EIF  | Low        |

### 2.2. Unadjusted Function Points (UFP)

Next, we use standard weights to calculate the UFP based on the number of components of each type and their complexity.

**Standard Weights Table:**
| Type | Low | Average | High |
| :--- | :-- | :------ | :--- |
| **EI** | 3   | 4       | 6    |
| **EO** | 4   | 5       | 7    |
| **EQ** | 3   | 4       | 6    |
| **ILF**| 7   | 10      | 15   |
| **EIF**| 5   | 7       | 10   |

**UFP Calculation Breakdown:**
*   **External Inputs (EI):**
    *   1 Average (User Profile) * 4 + 1 High (Route Inputs) * 6 = 4 + 6 = **10**
*   **External Outputs (EO):**
    *   1 Average (Display Route) * 5 + 2 Low (Error Msg, Export) * 4 = 5 + 8 = **13**
*   **External Inquiries (EQ):**
    *   1 Average (View History) * 4 = **4**
*   **Internal Logical Files (ILF):**
    *   2 Average (User Data, History Data) * 10 = **20**
*   **External Interface Files (EIF):**
    *   1 Low (Geo-data Service) * 7 = **7**

*   **Total UFP = 10 (EI) + 13 (EO) + 4 (EQ) + 20 (ILF) + 7 (EIF) = 54**

### 2.3. Value Adjustment Factor (VAF)

We assess 14 General System Characteristics (GSCs) on a scale of 0 (not present) to 5 (strong influence). The sum of these scores is the Total Degree of Influence (DI).

| #  | Characteristic | Score (0-5) |
| :- | :------------- | :---------- |
| 1  | Data Communications | 4 |
| 2  | Distributed Data Processing | 5 |
| 3  | Performance | 4 |
| 4  | Heavily Used Configuration | 3 |
| 5  | Transaction Rate | 3 |
| 6  | Online Data Entry | 4 |
| 7  | End-User Efficiency | 4 |
| 8  | Online Update | 3 |
| 9  | Complex Processing | 5 |
| 10 | Reusability | 3 |
| 11 | Installation Ease | 4 |
| 12 | Operational Ease | 4 |
| 13 | Multiple Sites | 5 |
| 14 | Facilitate Change | 3 |
| **Total Degree of Influence (DI):** | **53** |

**VAF Calculation:**
The VAF is calculated using the standard formula:
*   VAF = (DI * 0.01) + 0.65
*   VAF = (53 * 0.01) + 0.65 = 0.53 + 0.65 = **1.18**

### 2.4. Final Adjusted Function Points (FP)

The final FP is the product of the UFP and the VAF.
*   FP = UFP * VAF
*   FP = 54 * 1.18 = **63.72**

## 3. COCOMO II Estimation

The COCOMO II (Constructive Cost Model) is used to estimate software development effort and schedule.

### 3.1. Size in KLOC (Kilo Lines of Code)

First, we convert the abstract Function Points (FP) measurement into a physical size estimate (Source Lines of Code) using language-specific multipliers.

*   **Language Multipliers (LOC per FP):**
    *   JavaScript (Frontend): ~53
    *   Go (Backend): ~50
    *   Average LOC/FP = (53 + 50) / 2 = 51.5. We use **52** as a close approximation.

**Calculation:**
*   SLOC = FP * Average LOC/FP
*   SLOC = 63.72 * 52 = 3313.44 ≈ **3313**
*   Size (KLOC) = SLOC / 1000 = **3.313**

### 3.2. Effort Estimation

We use the **Semi-Detached** model, suitable for a project with a mix of experienced and inexperienced team members working on a familiar system.

*   **Formula:** Effort (E) in Person-Months = `a * (Size_in_KLOC)^b`
*   **Constants (Semi-Detached):**
    *   `a = 3.0`
    *   `b = 1.12`

**Calculation:**
*   E = 3.0 * (3.313)^1.12
*   E = 3.0 * 3.79 = **11.37 Person-Months**

### 3.3. Schedule Estimation

Finally, we estimate the total development time based on the calculated effort.

*   **Formula:** Time to Develop (TDEV) in Months = `c * (Effort)^d`
*   **Constants (Semi-Detached):**
    *   `c = 2.5`
    *   `d = 0.35`

**Calculation:**
*   TDEV = 2.5 * (12.2)^0.35
*   TDEV = 2.5 * 2.67 = 6.675 ≈ **6.7 Months**
