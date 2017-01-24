package time_and_distance

import (
	"fmt"
)

func	ComputeReverseRouteValues ( ) {

	Current := RightEndNode 
	ToRightEndNode := 0.0
	Distance := 0.0

	currentNode := Nodes [ RightEndNode ]

	currentNode.ToRightEndNode = ToRightEndNode
	currentNode.ToRightNode = Distance

	Nodes [ RightEndNode ] = currentNode

	fmt.Println ( "--------------------------------------" )

	for {
/*		Distance to Left Node ( aks, Previous Node ) is the same as 
		the Distance to Right Node for the Previous Nodeo
*/
		Previous := Nodes [ Current ].LeftNode

		fmt.Println ( "Reverse Walk:", Previous, "<--", Current )

		PreviousNode := Nodes [ Previous ]	// Previous Node

		Distance = Nodes [ Current ].ToLeftNode

//		fmt.Println ( Previous, "to Right :", PreviousNode.ToRightNode )
//		fmt.Println ( Previous, "to Right (New) :", Distance )

//		PreviousNode.ToRightNode = Distance

//		Runing total of Distance to Right End Node

		ToRightEndNode = ToRightEndNode + Distance

		fmt.Println ( "To Route End:", ToRightEndNode )

		PreviousNode.ToRightEndNode = ToRightEndNode

//		Do Time later

		Nodes [ Previous ] = PreviousNode

		fmt.Println ( Previous, ":", Nodes [ Previous ] )

	fmt.Println ( "--------------------------------------" )
		if Current == LeftEndNode { break }

		Current = Previous

	}
}



