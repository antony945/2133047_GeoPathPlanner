Groups and projects - Collection form and instructions

Here is the form for registering groups and projects : 

https://docs.google.com/forms/d/1zq_wUHi-1-wyrgOWS14ky7or3zXOtKgKzEFbxeCwGiI/preview

Registrations should happen within the specific cut-off date, in order to be able to discuss your project in the current academic year. In particular 15 January 2025 is the cut-off date for those ones willing to discuss projects before November 2025 (current roster). After the cut-off date, the instructor will provide feedbacks on the idea and will approve/ask for revision the project, within about one week (by the 20 January 2025). In case you do not receive such feedbacks, you can solicit the instructor.

Projects should be presented by groups of about 3/4 persons. 2 persons are possible as well, and up to 5 persons. The suggested dimension is 3/4 persons. For no reason they are possible projects by a single person or larger than 5 persons. Use the classroom to search for colleagues for the project, the classroom serves for this as well.

The initial idea (pitch) should be a text of a few paragraphs describing what are the objectives of the project, what the system should do, what are the potential users and the more important use cases. After the idea has been approved/revised, you will start specifying the user stories of the system. Suggestions for proposing a suitable project idea
The envisioned system should be distributed 
It should not be ONLY yet-another-WWW-application/site 
The document should not give the clear impression of very scarce creativity and attention in thinking and writing it (in many cases from previous years, it was self evident that students took 10 minutes to write in a hurry). The most of proposals are often simple, traditional, naive applications (booking of some items, rating apps, ecc.)
Conversely your specifications should: 
be creative, envisioning challenging systems
imply a distributed nature of the system
define clearly the functionalities you would like to offer
As a joke, keep this in mind: if your proposal is similar to a potential one for the Laboratorio di architetture software e sicurezza informatica, then it is not fine. You should challenge yourself more, and more creatively (also in order not to reuse previous projects).
Also to be more explicit : it is not sufficient to take a specification of a text of the course of Basi di Dati and/or Progettazione del Software and transform it into a WWW application. Again, too naive and too trivial.

Conducting and presenting the projects

Here are the instructions for conducting and presenting your projects
Register your group and idea / high-level specification within the cut-off date. The registration is binding, i.e., if not registered by the cut-off date, you cannot discuss your project. But then you can discuss it whenever you want during the year. 
After the approval of the idea by the instructor, you have to define the requirements by defining the user stories of your entire system. User stories can be documented by using spreadsheets (as the one proposed in a post), and collecting all of them at the end of the project in a booklet. For each user story, it should be provided also a LoFi mockup, to be prepared with the suggested tools (e.g., Balsamiq), and a textual description highlighting specific non functional requirements (if any).
After the definition of the user stories/requirements, you need to estimate the complexity of the software development, time and effort to carry it out, by using Function Points and COCOMO II methods. The analysis should be presented through spreadsheets (see the material) and a booklet explaining and detailing your method. At the end of the academic year, it will be interesting to compare estimation with actual KLOCs of the projects.
After requirements and estimation of the effort, you may design your project and document the software architecture. The development process should be based on SCRUM. It is required that you define the different sprints you may want to adopt for the development, each one with goals and planning. You can document the SCRUM method by adopting specific spreadsheets that have been published.
At the end, all the work done (design of the system, software architecture, sprints with relevant analytics, e.g., burndown charts) should be documented in a booklet.
The system is being developed in whatever technologies/framework you may want. The release should be done by providing the link of a GitHub repo with all the source code, configuration files, any other file you may need (please remember we have adopted a IaC approach - Infrastructure as Code) AND the Dockerfiles, docker-compose files, etc. which will allow the instructors to re-build/re-deploy your system on whatever platform (either on-premise or on cloud).
At the discussion of the project, you may need to take
a laptop with all the system running, to be used for a demo of the project
slides to be used (maximum 15 minutes of talk) for presenting the project idea, development, etc. 
all the booklets produced during the development process
Before the discussion, when booking your discussion, you will be required to provide the items below.

The discussion is booked by sending an email (the exact format of the email will be described in following posts, one for each session) and the instructor will take an appointment with the group for discussing the project. Please remind that all exams are IN-PRESENCE, so all group members will have to attend physically the discussion. Discussions will be held in a room with an HDMI projector, so please take a laptop for presenting your slides and your demo, and the adapter for your specific laptop in order to connect to the HDMI cable. The discussion of the project consists in a short talk with slides (maximum 15 minutes of talk) which presents the project idea, development, etc., and a demo of the system being built (via Docker tools) and running. The overall length of the discussion is approx. 25 minutes for each group.

Submission of the project material

A link to a GitHub repo, containing the following items.

1. The following textual items. All the textual documents must follow the Markdown syntax. Please refer to https://www.markdownguide.org/basic-syntax/ if you are not entirely familiar with it.
Following is the list of textual files to be submitted:
input.txt, which contains the description of the system (the one submitted in the form for approval) and the list of user stories.
Structure of the file: an example of input.txt, cf. https://drive.google.com/file/d/10mAO_d-HR4ubnqR8bcofkAYdddfDJKdu/view?usp=sharing 
Student_doc.md, which contains the specifics of the deployed system.
Structure of the file: Student_doc.md, cf. https://drive.google.com/file/d/1stCQoen6ojT3hBexAkyp0Ja8H6XzOuFn/view?usp=sharing
An example of a Student_doc.md file is accessible at https://drive.google.com/file/d/15lmOqwYTG4qORk3bZgwMlRBWLaA5Gn-I/view?usp=sharing
DataMetrics.json, which contains the features/qualities of the system. The file is created to categorize and group the user stories and analyze the dependencies and constraints they pose to the architectural design of the derived system. In particular, the user stories are grouped based on their real-world scope, i.e., user stories that have the same context and a similar scope are in the same group. Further, the creation of a group is based on the following factors: (a) user stories related to the same real-world object or situation are grouped together, and (b) user stories of different users who perform the same action are grouped together. From an architectural design point of view, all the user stories that belong to the same group have to be fulfilled by the same container. In this way, the system ensures the separation of concerns between containers while maintaining scope-related user stories grouping. How the user stories are fulfilled by the microservices inside the container is a design choice not inspected by this file.Structure of the file: DataMetrics.json, cf. https://drive.google.com/file/d/1_oJ8PKU7dzI4KA2uGuO-20JnJKX1Vj2A/view?usp=sharing
- set_id: numerical id of the set
- set_name: name of the set
- user_stories: list of ids of the user stories
- links: list of ids of other sets having related context
- db: need to store or retrieve data
From an architectural point of view, user stories that belong to linked sets can be fulfilled by the same container and the sets of user stories that are required to store or retrieve data must be fulfilled by a container hosting a database microservice.An example of DataMetrics.json file is accessible at https://drive.google.com/file/d/1p4nO8uYVAY9V-_LL5e2XG1_bQT00pRb3/view?usp=sharing. A sketch like https://drive.google.com/file/d/1oLtqBTFDqOM2x-r8IUq0YrCrt88b5CqK/view?usp=sharing can be helpful in constructing this file.
2. A folder with the following documents/booklets
The estimation of the complexity of the software development, time and effort to carry it out, by using Function Points and COCOMO II methods. The analysis should be presented through spreadsheets and a booklet explaining and detailing your method. 
The development process is based on SCRUM. It is required that you show the different sprints you may want to adopt for the development, each one with goals and planning. You can document the SCRUM method by adopting specific spreadsheets that have been published.
All the work done (design of the system, software architecture, sprints with relevant analytics, e.g., burndown charts) should be documented in a booklet.
3. A folder with the system developed in whatever technologies/framework you may want. In this folder all the source code, configuration files, and any other file you may need (please remember we have adopted an IaC approach - Infrastructure as Code) AND the Dockerfiles, docker-compose files, etc. which will allow the instructors to re-build/re-deploy your system on whatever platform (either on-premise or on cloud).

So, to wrap up, the structure of the repo submitted should be like this one, cf. https://drive.google.com/drive/folders/1HO9h8lQ8udyBhqjNk8BJ8hAODIONItcA?usp=sharing:
.
├── input.txt
├── Student_doc.md
├── DataMetrics.json
├── source    # directory containing the source files of the developed system
└── booklets # directory containing the documents about the development process, estimation, slides to be used for the discussion, etc.

Each group needs to create a repo named <MATRICOLA>_<PROJECT> where <MATRICOLA> is the INFOSTUD Student ID of the group leader and <PROJECT> is the acronym of your project.

In order to avoid misunderstandings, this is repeated.
Before the discussion, you need to share with the instructor all the documents and materials of your project, in particular
The requirements (as user stories) of your entire system. User stories can be documented by using a spreadsheet and collecting all of them in a booklet. For each user story, it should be provided also a LoFi mockup, to be prepared with the suggested tools (e.g., Balsamiq), and a textual description highlighting specific non functional requirements (if any).
The estimation of the complexity of the software development, time and effort to carry it out, by using Function Points and COCOMO II methods. The analysis should be presented through spreadsheets and a booklet explaining and detailing your method. 
The development process based on SCRUM. It is required that you show the different sprints you may want to adopt for the development, each one with goals and planning. You can document the SCRUM method by adopting specific spreadsheets that have been published.
All the work done (design of the system, software architecture, sprints with relevant analytics, e.g., burndown charts) should be documented in a booklet.
The system developed in whatever technologies/framework you may want. The release should be done by providing the link of a GitHub repo with all the source code, configuration files, any other file you may need (please remember we have adopted a IaC approach - Infrastructure as Code) AND the Dockerfiles, docker-compose files, etc. which will allow the instructors to re-build/re-deploy your system on whatever platform (either on-premise or on cloud).
Specific textual files (input.txt, DataMetrics.json, Student_doc.md)
All the above materials should be provided as a link to a GitHub repo with everything or it, or in any way you may deem appropriate.

If something is not clear, please ask and comment